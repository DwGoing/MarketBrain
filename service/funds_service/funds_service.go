package funds_service

import (
	context "context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	sync "sync"
	"time"

	"funds-system/pkg/bus_module"
	"funds-system/pkg/chain_module"
	"funds-system/pkg/hd_wallet"
	"funds-system/pkg/shared"
	"funds-system/pkg/storage_module"

	linq "github.com/ahmetb/go-linq/v3"
	"github.com/alibaba/ioc-golang/extension/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewFundsService
type FundsService struct {
	UnimplementedFundsServiceServer

	Mnemonic          *config.ConfigString          `config:",funds.mnemonic"`
	WalletMaxNumber   *config.ConfigInt64           `config:",funds.walletMaxNumber"`
	ExpireTime        *config.ConfigInt64           `config:",funds.expireTime"`
	ExpireDelay       *config.ConfigInt64           `config:",funds.expireDelay"`
	CollectThresholds *config.ConfigMap             `config:",funds.collectThresholds"`
	StorageModule     *storage_module.StorageModule `singleton:""`
	BusModule         *bus_module.BusModule         `singleton:""`
	ChainModule       *chain_module.ChainModule     `singleton:""`

	logger *log.Logger
}

/*
@title	构造函数
@param 	service 	*FundsService 	服务实例
@return _ 			*FundsService 	服务实例
@return _ 			error 			异常信息
*/
func NewFundsService(service *FundsService) (*FundsService, error) {
	service.logger = log.New(os.Stderr, "[FundsService]", log.LstdFlags)
	service.BusModule.RechargePaid = make(chan shared.RechargeRecord, 1024)
	// 开启充值监听
	go service.listenRecharge()
	return service, nil
}

/*
@title	充值监听
@param 	Self 	*FundsService 	服务实例
*/
func (Self *FundsService) listenRecharge() {
	confirmRecharge := func() {
		ctx := context.Background()
		wait := sync.WaitGroup{}
		redis, err := Self.StorageModule.GetRedisConnection()
		if err != nil {
			Self.logger.Printf("listenRecharge error: %s", err)
			return
		}
		defer redis.Close()
		mysql, err := Self.StorageModule.GetMysqlConnection()
		if err != nil {
			Self.logger.Printf("listenRecharge error: %s", err)
			return
		}
		sqlDB, err := mysql.DB()
		if err != nil {
			Self.logger.Printf("listenRecharge error: %s", err)
			return
		}
		defer sqlDB.Close()
		keys := redis.Keys(ctx, "RECHARGING_*").Val()
		for _, key := range keys {
			wait.Add(1)
			go func(key string) {
				defer wait.Done()
				// 获取充值信息
				recordId, err := redis.Get(ctx, key).Result()
				if err != nil {
					Self.logger.Printf("listenRecharge Error: %s", err)
					return
				}
				// 查询充值状态
				var record shared.RechargeRecord
				result := mysql.Table("RECHARGE_RECORD").Where("ID = ?", recordId).FirstOrInit(&record)
				if result.Error != nil {
					Self.logger.Printf("listenRecharge Error: %s", result.Error)
					return
				}
				if len(record.Id) <= 0 || record.Status == shared.RechargeStatus_PAID {
					// 移除钱包Recharging状态
					_, err = redis.Del(ctx, key).Result()
					if err != nil {
						Self.logger.Printf("listenRecharge Error: %s", err)
						return
					}
					return
				}
				// 检查余额
				balance, err := Self.ChainModule.GetBalance(record.WalletAddress, record.Token)
				if err != nil {
					Self.logger.Printf("listenRecharge Error: %s", err)
					return
				}
				beforeBalance, ok := new(big.Int).SetString(record.BeforeBalance, 10)
				if !ok {
					Self.logger.Println("listenRecharge Error: before balance invaild")
					return
				}
				amount, ok := new(big.Int).SetString(record.Amount, 10)
				if !ok {
					Self.logger.Println("listenRecharge Error: amount invaild")
					return
				}
				if new(big.Int).Add(beforeBalance, amount).Cmp(balance) > 0 {
					return
				}
				// 充值支付完成
				record.UpdatedAt = time.Now()
				record.AfterBalance = balance.String()
				record.Status = shared.RechargeStatus_PAID
				result = mysql.Table("RECHARGE_RECORD").Where("`ID` = ?", record.Id).Save(&record)
				if result.Error != nil {
					Self.logger.Printf("listenRecharge Error: %s", result.Error)
					return
				}
				// 通知消息总线
				Self.BusModule.RechargePaid <- record
				// 归集检查
				if threshold, ok := Self.CollectThresholds.Value()[record.Token]; ok {
					thresholdString := fmt.Sprint(threshold)
					go func(index uint32, token string) {
						thresholdBigFloat, ok := new(big.Float).SetString(thresholdString)
						if !ok {
							return
						}
						convertedThreshold, err := Self.ChainModule.ConvertValue(token, thresholdBigFloat)
						if err != nil {
							return
						}
						if balance.Cmp(convertedThreshold) >= 0 {
							hdWallet, err := Self.ChainModule.GetHDWallet(Self.Mnemonic.Value(), "")
							if err != nil {
								return
							}
							wallet, err := hdWallet.GetWallet(index)
							if err != nil {
								return
							}
							err = Self.collect(wallet, common.HexToAddress("0x0"), common.HexToAddress(token))
							if err != nil {
								return
							}
						}
						log.Println("归集成功")
					}(uint32(record.WalletIndex), record.Token)
				}
				// 移除钱包Recharging状态
				_, err = redis.Del(ctx, key).Result()
				if err != nil {
					Self.logger.Printf("listenRecharge Error: %s", err)
					return
				}
			}(key)
		}
		wait.Wait()
	}
	for {
		confirmRecharge()
		time.Sleep(time.Millisecond * 500)
	}
}

/*
@title	转账
@param 	Self 	*FundsService 		服务实例
@param 	from 	*hd_wallet.Wallet 	发送钱包
@param 	to 		common.Address 		接收地址
@param 	token 	common.Address 		转账Token
@param 	amount 	*big.Int 			数量
@param 	remarks string 				备注
@return _ 		error 				异常信息
*/
func (Self *FundsService) transfer(from *hd_wallet.Wallet, to common.Address, token common.Address, amount *big.Int, remarks string) error {
	if from == nil {
		return errors.New("from nil")
	}
	mysql, err := Self.StorageModule.GetMysqlConnection()
	if err != nil {
		return err
	}
	sqlDB, err := mysql.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	record := shared.TransferRecord{
		Record: shared.Record{
			Id:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		FromIndex:   from.Index,
		FromAddress: from.GetAddress(),
		To:          to.Hex(),
		Token:       token.Hex(),
		Amount:      amount.String(),
	}
	// 创建转账记录
	result := mysql.Table("TRANSFER_RECORD").Create(&record)
	if result.Error != nil {
		return result.Error
	}
	// 执行转账操作
	err = Self.ChainModule.Transfer(from.PrivateKey, to.Hex(), token.Hex(), amount)
	// 更新转账状态
	if err != nil {
		record.Status = shared.TransferStatus_FAILED
		record.Error = err.Error()
	} else {
		record.Status = shared.TransferStatus_SUCCESS
	}
	record.UpdatedAt = time.Now()
	result = mysql.Table("TRANSFER_RECORD").Where("`ID`=?", record.Id).Save(&record)
	if result.Error != nil {
		return result.Error
	}
	return err
}

/*
@title	资金归集
@param 	Self 	*FundsService 		服务实例
@param 	from 	*hd_wallet.Wallet 	发送钱包
@param 	to 		common.Address 		接收地址
@param 	token 	common.Address 		转账Token
@return _ 		error 				异常信息
*/
func (Self *FundsService) collect(from *hd_wallet.Wallet, to common.Address, token common.Address) error {
	hdWallet, err := Self.ChainModule.GetHDWallet(Self.Mnemonic.Value(), "")
	if err != nil {
		return err
	}
	// 检查ETH余额
	ethBalance, err := Self.ChainModule.GetBalance(from.GetAddress(), "0x0")
	if err != nil {
		return err
	}
	// 计算转账数量
	gasPrice, err := Self.ChainModule.GetGasPrice()
	if err != nil {
		return err
	}
	var gas *big.Int
	var amount *big.Int
	if token == common.HexToAddress("0x0") {
		// ETH转账
		gas = new(big.Int).Mul(big.NewInt(21000), gasPrice)
		if ethBalance.Cmp(gas) <= 0 {
			return errors.New("balance less then gas")
		}
		amount = new(big.Int).Sub(ethBalance, gas)
	} else {
		// ERC20转账
		balance, err := Self.ChainModule.GetBalance(from.GetAddress(), token.Hex())
		if err != nil {
			return err
		}
		if balance.Cmp(big.NewInt(0)) <= 0 {
			return errors.New("balance less then zero")
		}
		gas = new(big.Int).Mul(big.NewInt(300000), gasPrice)
		// Gas不足
		if ethBalance.Cmp(gas) <= 0 {
			gasWallet, err := hdWallet.GetWallet(0)
			if err != nil {
				return err
			}
			err = Self.transfer(
				gasWallet,
				common.HexToAddress(from.GetAddress()),
				common.HexToAddress("0x0"),
				new(big.Int).Mul(gas, big.NewInt(10)), // 补充10倍Gas
				"补充Gas",
			)
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * 30000) // 30秒后重新发起归集
			return Self.collect(from, to, token)
		}
		amount = balance
	}
	if to == common.HexToAddress("0x0") {
		// 默认发送到#0钱包
		wallet, err := hdWallet.GetWallet(0)
		if err != nil {
			return err
		}
		to = common.HexToAddress(wallet.GetAddress())
	}
	err = Self.transfer(from, to, token, amount, "资金归集")
	if err != nil {
		return err
	}
	return nil
}

/*
@title			获取归集钱包
@description	归集钱包为#0钱包
@param 			Self	*FundsService 		 		服务实例
@param 			ctx 	context.Context 			请求上下文
@param 			request	*emptypb.Empty 				请求体
@return 		_ 		*GetWalletAddressResponse 	响应体
@return 		_ 		error 						异常信息
*/
func (Self *FundsService) GetCollectionWallet(ctx context.Context, request *emptypb.Empty) (*GetCollectionWalletResponse, error) {
	hdWallet, err := Self.ChainModule.GetHDWallet(Self.Mnemonic.Value(), "")
	if err != nil {
		return nil, err
	}
	wallet, err := hdWallet.GetWallet(0)
	if err != nil {
		return nil, err
	}
	address := wallet.GetAddress()
	balance, err := Self.ChainModule.GetBalance(address, "0x0")
	if err != nil {
		return nil, err
	}
	return &GetCollectionWalletResponse{
		Address: address,
		Balance: balance.Uint64(),
	}, nil
}

/*
@title			获取充值钱包
@description	充值钱包默认为非#0钱包，每个钱包有效期为expireTime + expireDelay
@param 			Self 	*FundsService 		 		服务实例
@param 			ctx 	context.Context 			请求上下文
@param 			request *GetWalletAddressRequest 	请求体
@return 		_ 		*GetWalletAddressResponse 	响应体
@return 		_ 		error 						异常信息
*/
func (Self *FundsService) GetRechargeWallet(ctx context.Context, request *GetRechargeWalletRequest) (*GetRechargeWalletResponse, error) {
	// 检查参数
	if len(request.ExternalIdentity) <= 0 {
		return nil, errors.New("GetRechargeWallet Error: exteranl identity empty")
	}
	if len(request.CallbackUrl) <= 0 {
		return nil, errors.New("GetRechargeWallet Error: callback url empty")
	}
	if len(request.Amount) <= 0 {
		return nil, errors.New("GetRechargeWallet Error: amount empty")
	}
	mysql, err := Self.StorageModule.GetMysqlConnection()
	if err != nil {
		return nil, err
	}
	sqlDB, err := mysql.DB()
	if err != nil {
		Self.logger.Printf("GetRechargeWallet error: %s", err)
		return nil, err
	}
	defer sqlDB.Close()
	var count int64
	result := mysql.Table("RECHARGE_RECORD").Where("EXTERNAL_IDENTITY = ?", request.ExternalIdentity).Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}
	if count > 0 {
		return nil, errors.New("GetRechargeWallet Error: external identity existed")
	}
	amount, ok := new(big.Float).SetString(request.Amount)
	if !ok {
		err = errors.New("amount invaild")
		Self.logger.Printf("GetRechargeWallet error: %s", err)
		return nil, err
	}
	convertedAmount, err := Self.ChainModule.ConvertValue(request.Token, amount)
	if err != nil {
		return nil, err
	}
	hdWallet, err := Self.ChainModule.GetHDWallet(Self.Mnemonic.Value(), "")
	if err != nil {
		return nil, err
	}
	redis, err := Self.StorageModule.GetRedisConnection()
	if err != nil {
		return nil, err
	}
	defer redis.Close()
	var record shared.RechargeRecord
	retry := 0
	for {
		// 重试10次
		if retry++; retry > 10 {
			return nil, errors.New("GetRechargeWallet Error: no enough available wallet")
		}
		// 获取一个钱包索引（原子操作）
		index, err := redis.Eval(ctx,
			`local key = KEYS[1]
			local value = redis.call("GET", key)
			if not value or tonumber(value) >= tonumber(ARGV[1]) then
				value = 0
			end
			value = value + 1
			redis.call("SET", key, value)
			return value`,
			[]string{"WALLET_INDEX"}, Self.WalletMaxNumber.Value(),
		).Result()
		if err != nil {
			Self.logger.Printf("GetRechargeWallet Error: %s", err)
			continue
		}
		wallet, err := hdWallet.GetWallet(uint32(index.(int64)))
		if err != nil {
			Self.logger.Printf("GetRechargeWallet Error: %s", err)
			continue
		}
		// 检查Recharging状态
		if redis.Exists(ctx, "RECHARGING_"+wallet.GetAddress()).Val() > 0 {
			Self.logger.Printf("GetRechargeWallet Error: status is recharging")
			continue
		}
		// 检查Collecting状态
		if redis.Exists(ctx, "COLLECTING_"+wallet.GetAddress()).Val() > 0 {
			Self.logger.Printf("GetRechargeWallet Error: status is collecting")
			continue
		}
		// 检查原始余额
		balance, err := Self.ChainModule.GetBalance(wallet.GetAddress(), request.Token)
		if err != nil {
			Self.logger.Printf("GetRechargeWallet Error: %s", err)
			continue
		}
		duration := time.Second * time.Duration(Self.ExpireTime.Value()+Self.ExpireDelay.Value())
		record = shared.RechargeRecord{
			Record: shared.Record{
				Id:        uuid.NewString(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ExternalIdentity: request.ExternalIdentity,
			ExternalData:     request.ExternalData,
			CallbackUrl:      request.CallbackUrl,
			Token:            common.HexToAddress(request.Token).Hex(),
			Amount:           convertedAmount.String(),
			WalletIndex:      index.(int64),
			WalletAddress:    wallet.GetAddress(),
			BeforeBalance:    balance.String(),
			Status:           shared.RechargeStatus_UNPAID,
			ExpireAt:         time.Now().Add(duration),
		}
		// 添加钱包Recharging状态
		// expireTime + expireDelay
		ok, err := redis.SetNX(ctx, "RECHARGING_"+wallet.GetAddress(), record.Id, duration).Result()
		if !ok {
			return nil, err
		}
		result = mysql.Table("RECHARGE_RECORD").Create(&record)
		if result.Error != nil {
			return nil, result.Error
		}
		break
	}
	return &GetRechargeWalletResponse{
		Id:       record.Id,
		Address:  record.WalletAddress,
		ExpireAt: record.ExpireAt.Unix(),
	}, nil
}

/*
@title	获取充值记录
@param 	Self 	*FundsService 				服务实例
@param 	ctx 	context.Context 			请求上下文
@param 	request *GetRechargeRecordsRequest 	请求体
@return _ 		*GetRechargeRecordsResponse 响应体
@return _ 		error 						异常信息
*/
func (Self *FundsService) GetRechargeRecords(ctx context.Context, request *GetRechargeRecordsRequest) (*GetRechargeRecordsResponse, error) {
	// 处理参数
	pageSize := int(request.PageSize)
	if pageSize == 0 {
		pageSize = 20 // 默认一页20条记录
	}
	pageIndex := int(request.PageIndex)
	if pageIndex == 0 {
		pageIndex = 1
	}
	mysql, err := Self.StorageModule.GetMysqlConnection()
	if err != nil {
		return nil, err
	}
	sqlDB, err := mysql.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()
	db := mysql.Table("RECHARGE_RECORD")
	if len(request.Conditions) > 0 {
		db = db.Where(request.Conditions)
	}
	var total int64
	result := db.Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}
	var rechargeRecords []shared.RechargeRecord
	result = db.Order("CREATED_AT DESC").Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&rechargeRecords)
	if result.Error != nil {
		return nil, result.Error
	}
	var rechargeOrders []*RechargeRecord
	linq.From(rechargeRecords).SelectT(func(item shared.RechargeRecord) *RechargeRecord {
		amount, ok := new(big.Int).SetString(item.Amount, 10)
		if !ok {
			return nil
		}
		unconvertedValue, err := Self.ChainModule.UnconvertValue(item.Token, amount)
		if err != nil {
			return nil
		}
		return &RechargeRecord{
			Id:               item.Id,
			CreatedAt:        item.CreatedAt.String(),
			UpdatedAt:        item.UpdatedAt.String(),
			ExternalIdentity: item.ExternalIdentity,
			ExternalData:     item.ExternalData,
			CallbackUrl:      item.CallbackUrl,
			Token:            item.Token,
			Amount:           unconvertedValue.Text('f', 6),
			WalletAddress:    item.WalletAddress,
			Status:           item.Status.String(),
			ExpireAt:         item.ExpireAt.String(),
		}
	}).ToSlice(&rechargeOrders)
	return &GetRechargeRecordsResponse{
		Result: rechargeOrders,
		Total:  total,
	}, nil
}

// @title	资金归集结果
type FundsCollectResult struct {
	WalletIndex   uint32
	WalletAddress string
	Error         error
}

/*
@title			资金归集
@description	归集资金到#0钱包或指定钱包
@param 			Self	*FundsService 			服务实例
@param 			ctx 	context.Context 		请求上下文
@param 			request	*FundsCollectRequest 	请求体
@return 		_ 		*emptypb.Empty		 	响应体
@return 		_ 		error 					异常信息
*/
func (Self *FundsService) FundsCollect(ctx context.Context, request *FundsCollectRequest) (*emptypb.Empty, error) {
	hdWallet, err := Self.ChainModule.GetHDWallet(Self.Mnemonic.Value(), "")
	if err != nil {
		return nil, err
	}
	channel := make(chan FundsCollectResult, 1024)
	// 任务监听
	go func(c chan FundsCollectResult) {
		taskNumber := Self.WalletMaxNumber.Value()
		var current int64
		for {
			select {
			case result := <-c:
				current++
				log.Printf("%s %s", result.WalletAddress, result.Error)
				if result.Error == nil {
					continue
				}
			default:
				if current >= taskNumber {
					return
				}
				time.Sleep(time.Millisecond * 500)
			}
		}
	}(channel)
	// 转账任务
	for i := 1; i <= int(Self.WalletMaxNumber.Value()); i++ {
		go func(index uint32, c chan FundsCollectResult) {
			result := FundsCollectResult{}
			wallet, err := hdWallet.GetWallet(index)
			if err != nil {
				result.Error = err
				c <- result
				return
			}
			result.WalletIndex = wallet.Index
			result.WalletAddress = wallet.GetAddress()
			err = Self.collect(
				wallet,
				common.HexToAddress(request.To),
				common.HexToAddress(request.Token),
			)
			if err != nil {
				result.Error = err
				c <- result
				return
			}
			// 转账成功
			c <- result
		}(uint32(i), channel)
		time.Sleep(time.Millisecond * 500)
	}
	return &emptypb.Empty{}, nil
}
