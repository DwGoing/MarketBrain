package funds_service

import (
	"bytes"
	context "context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	sync "sync"
	"time"

	"github.com/DwGoing/OnlyPay/docs"
	"github.com/DwGoing/OnlyPay/internal/module/chain_module"
	"github.com/DwGoing/OnlyPay/internal/module/config_module"
	"github.com/DwGoing/OnlyPay/internal/module/storage_module"
	"github.com/DwGoing/OnlyPay/internal/shared"
	"github.com/DwGoing/OnlyPay/pkg/hd_wallet"
	"github.com/ahmetb/go-linq"
	"github.com/alibaba/ioc-golang/extension/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewFundsService
type FundsService struct {
	GinPort       *config.ConfigInt64           `config:",app.gin.port"`
	GrpcPort      *config.ConfigInt64           `config:",app.grpc.port"`
	Nodes         *config.ConfigSlice           `config:",chain.nodes"`
	ConfigModule  *config_module.ConfigModule   `singleton:""`
	StorageModule *storage_module.StorageModule `singleton:""`
	ChainModule   *chain_module.ChainModule     `singleton:""`

	UnimplementedFundsServiceServer
	logger       *log.Logger
	RechargePaid chan shared.RechargeRecord
}

/*
@title	构造函数
@param 	service *FundsService 	服务实例
@return _ 		*FundsService 	服务实例
@return _ 		error 			异常信息
*/
func NewFundsService(service *FundsService) (*FundsService, error) {
	err := service.ChainModule.Initialize(service.Nodes.Value())
	if err != nil {
		return nil, err
	}
	service.logger = log.New(os.Stdout, "[FundsService]", log.LstdFlags)
	service.RechargePaid = make(chan shared.RechargeRecord, 1024)
	// 开启事件监听
	go service.listenEvent()
	// 开启充值监听
	go service.listenRecharge()
	return service, nil
}

/*
@title	初始化
@param 	Self	*FundsService 	服务实例
@return _ 		*FundsService 	服务实例
@return _ 		error 			异常信息
*/
func (Self *FundsService) Initialize() error {
	// gin
	go func() {
		docs.SwaggerInfo.BasePath = "/"
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Self.GinPort.Value()))
		if err != nil {
			Self.logger.Fatalf("Gin初始化失败:%s", err)
		}
		engine := gin.Default()
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		engine.POST("/callback", func(ctx *gin.Context) {
			body, _ := io.ReadAll(ctx.Request.Body)
			log.Printf("%s", body)
			defer ctx.Request.Body.Close()
		})
		v1Router := engine.Group("/v1")
		{
			configRouter := v1Router.Group("config")
			{
				configRouter.GET("/load", LoadConfig)
				configRouter.POST("/set", SetConfig)
			}
			fundsRouter := v1Router.Group("/funds")
			{
				fundsRouter.POST("/getRechargeWallet", GetRechargeWallet)
				fundsRouter.GET("/getRechargeRecords", GetRechargeRecords)
			}
		}
		Self.logger.Printf("Gin正在监听: %s", listener.Addr())
		if err := engine.RunListener(listener); err != nil {
			Self.logger.Fatalf("Gin启动失败: %s", err)
		}
	}()

	// gRPC
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Self.GrpcPort.Value()))
		if err != nil {
			Self.logger.Fatalf("gRPC初始化失败:%s", err)
		}
		server := grpc.NewServer()
		RegisterFundsServiceServer(server, Self)
		Self.logger.Printf("gRPC正在监听: %s", listener.Addr())
		if err = server.Serve(listener); err != nil {
			Self.logger.Fatalf("gRPC启动失败: %s", err)
		}
	}()
	return nil
}

/*
@title	监听事件
@param 	Self 	*FundsService 	服务实例
*/
func (Self *FundsService) listenEvent() {
	for {
		select {
		// 充值完成
		case record := <-Self.RechargePaid:
			Self.rechargePaidHandle(record)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

/*
@title	充值完成事件
@param 	Self 	*FundsService 			服务实例
@param 	record 	shared.RechargeRecord 	充值记录
*/
func (Self *FundsService) rechargePaidHandle(record shared.RechargeRecord) {
	log.Printf("%s 充值完成", record.Id)
	go func() {
		retry := 0
		for {
			// 重试5次
			if retry++; retry > 5 {
				log.Printf("rechargePaidHandle Error: maximum retry limit")
				return
			}
			time.Sleep(time.Minute * 2 * time.Duration(retry-1)) // 0/2/4/6/8 min
			request, err := http.NewRequest("POST", record.CallbackUrl, bytes.NewBuffer(record.ExternalData))
			if err != nil {
				log.Printf("rechargePaidHandle Error: %s", err)
				continue
			}
			httpResponse, err := http.DefaultClient.Do(request)
			if err != nil {
				log.Printf("rechargePaidHandle Error: %s", err)
				continue
			}
			defer httpResponse.Body.Close()
			if httpResponse.StatusCode != http.StatusOK {
				log.Printf("rechargePaidHandle Error: status code not 200")
				continue
			}
			return
		}
	}()
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
				// 通知事件总线
				Self.RechargePaid <- record
				// 归集检查
				configs, err := Self.ConfigModule.Load()
				if err != nil {
					return
				}
				if threshold, ok := configs.CollectThresholds[record.Token]; ok {
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
							hdWallet, err := Self.ChainModule.GetHDWallet(configs.Mnemonic, "")
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
	configs, err := Self.ConfigModule.Load()
	if err != nil {
		return err
	}
	hdWallet, err := Self.ChainModule.GetHDWallet(configs.Mnemonic, "")
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
@title	加载配置
@param 	Self	*FundsService 	服务实例
@param 	ctx		context.Context 请求上下文
@param 	request	*LoadRequest 	请求体
@return _ 		*emptypb.Empty 	响应体
@return _ 		error 			异常信息
*/
func (Self *FundsService) LoadConfig(ctx context.Context, request *emptypb.Empty) (*LoadConfigResponse, error) {
	configs, err := Self.ConfigModule.Load()
	if err != nil {
		return nil, err
	}
	return &LoadConfigResponse{
		Mnemonic:          configs.Mnemonic,
		WalletMaxNumber:   configs.WalletMaxNumber,
		ExpireTime:        configs.ExpireTime,
		ExpireDelay:       configs.ExpireDelay,
		CollectThresholds: configs.CollectThresholds,
	}, nil
}

/*
@title	修改配置
@param 	Self	*FundsService 	服务实例
@param 	ctx		context.Context 请求上下文
@param 	request	*SetRequest 	请求体
@return _ 		*emptypb.Empty 	响应体
@return _ 		error 			异常信息
*/
func (Self *FundsService) SetConfig(ctx context.Context, request *SetConfigRequest) (*emptypb.Empty, error) {
	if request.Mnemonic != nil {
		err := Self.ConfigModule.Set("Mnemonic", *request.Mnemonic)
		if err != nil {
			return nil, err
		}
	}
	if request.WalletMaxNumber != nil {
		err := Self.ConfigModule.Set("WalletMaxNumber", *request.WalletMaxNumber)
		if err != nil {
			return nil, err
		}
	}
	if request.ExpireTime != nil {
		err := Self.ConfigModule.Set("ExpireTime", *request.ExpireTime)
		if err != nil {
			return nil, err
		}
	}
	if request.ExpireDelay != nil {
		err := Self.ConfigModule.Set("ExpireDelay", *request.ExpireDelay)
		if err != nil {
			return nil, err
		}
	}
	if request.CollectThresholds != nil {
		err := Self.ConfigModule.Set("CollectThresholds", request.CollectThresholds)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
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
	configs, err := Self.ConfigModule.Load()
	if err != nil {
		return nil, err
	}
	hdWallet, err := Self.ChainModule.GetHDWallet(configs.Mnemonic, "")
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
	unconvertedBalance, err := Self.ChainModule.UnconvertValue("0x0", balance)
	if err != nil {
		return nil, err
	}
	return &GetCollectionWalletResponse{
		Address: address,
		Balance: unconvertedBalance.String(),
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
	configs, err := Self.ConfigModule.Load()
	if err != nil {
		return nil, err
	}
	hdWallet, err := Self.ChainModule.GetHDWallet(configs.Mnemonic, "")
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
			[]string{"WALLET_INDEX"}, configs.WalletMaxNumber,
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
		duration := time.Second * time.Duration(configs.ExpireTime+configs.ExpireDelay)
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
	configs, err := Self.ConfigModule.Load()
	if err != nil {
		return nil, err
	}
	hdWallet, err := Self.ChainModule.GetHDWallet(configs.Mnemonic, "")
	if err != nil {
		return nil, err
	}
	channel := make(chan FundsCollectResult, 1024)
	// 任务监听
	go func(c chan FundsCollectResult) {
		taskNumber := configs.WalletMaxNumber
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
	for i := 1; i <= int(configs.WalletMaxNumber); i++ {
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

/*
@title	从归集钱包转账
@param 	Self 	*FundsService 		服务实例
@param 	ctx 	context.Context 	请求上下文
@param 	request *TransferRequest 	接收地址
@return _ 		error 				异常信息
*/
func (Self *FundsService) Transfer(ctx context.Context, request *TransferRequest) (*emptypb.Empty, error) {
	configs, err := Self.ConfigModule.Load()
	if err != nil {
		return nil, err
	}
	hdWallet, err := Self.ChainModule.GetHDWallet(configs.Mnemonic, "")
	if err != nil {
		return nil, err
	}
	wallet, err := hdWallet.GetWallet(0)
	if err != nil {
		return nil, err
	}
	amount, ok := new(big.Float).SetString(request.Amount)
	if !ok {
		return nil, errors.New("amount invaild")
	}
	convertedAmount, err := Self.ChainModule.ConvertValue(request.Token, amount)
	if err != nil {
		return nil, err
	}
	err = Self.transfer(wallet, common.HexToAddress(request.To), common.HexToAddress(request.Token), convertedAmount, request.Remarks)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
