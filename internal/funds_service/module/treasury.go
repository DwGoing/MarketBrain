package module

import (
	context "context"
	"errors"
	"strings"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/internal/funds_service/module/treasury_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/DwGoing/MarketBrain/pkg/hd_wallet"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Treasury struct {
	treasury_generated.UnimplementedTreasuryServer
}

// @title	创建充值订单
// @param	Self				*Treasury		模块实例
// @param	externalIdentity	string			扩展标识
// @param	externalData		[]byte			扩展数据
// @param	callbackUrl			string			回调链接
// @param	chainType			enum.ChainType	主链类型
// @param	amount				float64			充值数量
// @param	walletIndex			int64			钱包索引
// @return	_					string			订单ID
// @return	_					string			接收钱包
// @return	_					error			异常信息
func (Self *Treasury) createRechargeOrder(
	externalIdentity string,
	externalData []byte,
	callbackUrl string,
	chainType string,
	amount float64,
	walletIndex int64,
) (string, string, time.Time, error) {
	if strings.TrimSpace(externalIdentity) == "" ||
		strings.TrimSpace(callbackUrl) == "" ||
		strings.TrimSpace(chainType) == "" ||
		amount < 1 ||
		walletIndex < 1 {
		return "", "", time.Time{}, errors.New("parameter invaild")
	}
	chain, err := new(enum.ChainType).Parse(chainType)
	if err != nil {
		return "", "", time.Time{}, err
	}
	configModule, err := GetConfig()
	if err != nil {
		return "", "", time.Time{}, err
	}
	config, err := configModule.Load()
	if err != nil {
		return "", "", time.Time{}, err
	}
	hdWallet, err := hd_wallet.FromMnemonic(config.Mnemonic, "")
	if err != nil {
		return "", "", time.Time{}, err
	}
	var wallet string
	switch chain {
	case enum.ChainType_TRON:
		account, err := hdWallet.GetAccount(hd_wallet.Currency_TRON, walletIndex)
		if err != nil {
			return "", "", time.Time{}, err
		}
		wallet = account.GetAddress()
	default:
		return "", "", time.Time{}, errors.New("unsupported chain")
	}
	storageModule, err := GetStorage()
	if err != nil {
		return "", "", time.Time{}, err
	}
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return "", "", time.Time{}, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return "", "", time.Time{}, err
	}
	defer db.Close()
	record := &model.RechargeOrderRecord{
		ExternalIdentity: externalIdentity,
		ExternalData:     externalData,
		CallbackUrl:      callbackUrl,
		ChainType:        chainType,
		Amount:           amount,
		Wallet:           wallet,
		ExpireAt:         time.Now().Add(time.Minute * 30),
	}
	record, err = model.CreateRechargeOrderRecord(mysqlClient, record)
	if err != nil {
		return "", wallet, time.Time{}, err
	}
	return record.Id, record.Wallet, record.ExpireAt, nil
}

// @title	创建充值订单
// @param	Self		*Treasury										模块实例
// @param	ctx			context.Context									上下文
// @param	request		*treasury_generated.CreateRechargeOrderRequest	请求体
// @return	_			*treasury_generated.CreateRechargeOrderResponse	响应体
// @return	_			error											异常信息
func (Self *Treasury) CreateRechargeOrderRpc(ctx context.Context, request *treasury_generated.CreateRechargeOrderRequest) (*treasury_generated.CreateRechargeOrderResponse, error) {
	orderId, wallet, expireAt, err := Self.createRechargeOrder(
		request.ExternalIdentity,
		request.ExternalData,
		request.CallbackUrl,
		request.ChainType,
		request.Amount,
		request.WalletIndex,
	)
	if err != nil {
		return nil, err
	}
	return &treasury_generated.CreateRechargeOrderResponse{
		OrderId:  orderId,
		Wallet:   wallet,
		ExpireAt: expireAt.String(),
	}, nil
}

type CreateRechargeOrderRequest struct {
	ExternalIdentity string  `json:"externalIdentity"`
	ExternalData     []byte  `json:"externalData"`
	CallbackUrl      string  `json:"callbackUrl"`
	ChainType        string  `json:"chainType"`
	Amount           float64 `json:"amount"`
	WalletIndex      int64   `json:"walletIndex"`
}

type CreateRechargeOrderResponse struct {
	OrderId  string    `json:"orderId"`
	Wallet   string    `json:"wallet"`
	ExpireAt time.Time `json:"expireAt"`
}

// @title	创建充值订单
// @param	Self	*Treasury		模块实例
// @param	ctx		*gin.Context	上下文
func (Self *Treasury) CreateRechargeOrderApi(ctx *gin.Context) {
	var request CreateRechargeOrderRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_RequestBindError, err)
	}
	orderId, wallet, expireAt, err := Self.createRechargeOrder(
		request.ExternalIdentity,
		request.ExternalData,
		request.CallbackUrl,
		request.ChainType,
		request.Amount,
		request.WalletIndex,
	)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
	}
	Response.Success(ctx, CreateRechargeOrderResponse{
		OrderId:  orderId,
		Wallet:   wallet,
		ExpireAt: expireAt,
	})
}

// @title	提交充值订单交易Hash
// @param	Self		*Treasury	模块实例
// @param	orderId		string		订单ID
// @param	txHash		string		交易Hash
func (Self *Treasury) submitRechargeOrderTransaction(orderId string, txHash string) error {
	if strings.TrimSpace(orderId) == "" ||
		strings.TrimSpace(txHash) == "" {
		return errors.New("parameter invaild")
	}
	storageModule, err := GetStorage()
	if err != nil {
		return err
	}
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return err
	}
	defer db.Close()
	// 检查订单
	rechargeOrders, total, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`ID` = ?",
		ConditionsParameters: []any{orderId},
	})
	if err != nil {
		return err
	}
	if total < 1 {
		return errors.New("order not existed")
	}
	// 检查TxHash
	_, total, err = model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`TX_HASH` = ?",
		ConditionsParameters: []any{txHash},
	})
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("tx hash existed")
	}
	rechargeOrder := rechargeOrders[0]
	status, _ := new(enum.RechargeStatus).Parse(rechargeOrder.Status)
	if status == enum.RechargeStatus_PAID {
		return nil
	}
	// 更新交易Hash
	err = model.UpdateRechargeOrderRecords(mysqlClient, model.UpdateOption{
		Conditions:           "`ID` = ?",
		ConditionsParameters: []any{orderId},
		Values: map[string]any{
			"TX_HASH": txHash,
			"STATUS":  enum.RechargeStatus_UNPAID.String(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// @title	提交充值订单交易Hash
// @param	Self	*Treasury													模块实例
// @param	ctx		context.Context												上下文
// @param	request	*treasury_generated.SubmitRechargeOrderTransactionRequest	请求体
// @return	_		*emptypb.Empty												响应体
// @return	_		error														异常信息
func (Self *Treasury) SubmitRechargeOrderTransactionRpc(ctx context.Context, request *treasury_generated.SubmitRechargeOrderTransactionRequest) (*emptypb.Empty, error) {
	err := Self.submitRechargeOrderTransaction(request.OrderId, request.TxHash)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

type SubmitRechargeOrderTransactionRequest struct {
	OrderId string `json:"orderId"`
	TxHash  string `json:"txHash"`
}

// @title	提交充值订单交易Hash
// @param	Self	*Treasury		模块实例
// @param	ctx		*gin.Context	上下文
// @return	_		error			异常信息
func (Self *Treasury) SubmitRechargeOrderTransactionApi(ctx *gin.Context) {
	var request SubmitRechargeOrderTransactionRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_RequestBindError, err)
		return
	}
	err = Self.submitRechargeOrderTransaction(request.OrderId, request.TxHash)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
		return
	}
	Response.Success(ctx, nil)
}

// @title	检查充值订单状态
// @param	Self	*Treasury	模块实例
// @return	_		error		异常信息
func (Self *Treasury) CheckRechargeOrderStatus() error {
	checkExpireTime := func(client *gorm.DB, order model.RechargeOrderRecord) {
		// 检查过期时间
		if order.ExpireAt.Before(time.Now()) {
			model.UpdateRechargeOrderRecords(client, model.UpdateOption{
				Conditions:           "`ID` = ?",
				ConditionsParameters: []any{order.Id},
				Values: map[string]any{
					"STATUS": enum.RechargeStatus_CANCELLED.String(),
				},
			})
		}
	}
	storageModule, err := GetStorage()
	if err != nil {
		return err
	}
	redisClient, err := storageModule.GetRedisClient()
	if err != nil {
		return err
	}
	defer redisClient.Close()
	// 加锁
	lock := "RECHARGE_ORDER_STATUS_CHEAKING"
	ok, err := redisClient.SetNX(context.Background(), lock, "", time.Duration(time.Minute*10)).Result()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	// 解锁
	defer redisClient.Del(context.Background(), lock).Result()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return err
	}
	defer db.Close()
	// 检查所有未支付订单
	rechargeOrders, _, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`STATUS` = ?",
		ConditionsParameters: []any{enum.RechargeStatus_UNPAID.String()},
	})
	if err != nil {
		return err
	}
	configModule, err := GetConfig()
	if err != nil {
		return err
	}
	config, err := configModule.load()
	if err != nil {
		return err
	}
	chain, err := GetChain()
	if err != nil {
		return err
	}
	for _, rechargeOrder := range rechargeOrders {
		chainType, _ := new(enum.ChainType).Parse(rechargeOrder.ChainType)
		chainConfig, ok := config.ChainConfigs[chainType.String()]
		if !ok {
			continue
		}
		// 检查交易状态
		if strings.TrimSpace(rechargeOrder.TxHash) != "" {
			result, address, timeStamp, to, amount, confirms, err := chain.DecodeTransaction(chainType, rechargeOrder.TxHash)
			if err != nil {
				zap.S().Errorf("decode transaction error: %s", err)
				// 检查是否过期
				checkExpireTime(mysqlClient, rechargeOrder)
				continue
			}
			if !result ||
				address != chainConfig.USDT ||
				timeStamp < rechargeOrder.CreatedAt.UTC().UnixMilli() ||
				to != rechargeOrder.Wallet ||
				amount < rechargeOrder.Amount {
				model.UpdateRechargeOrderRecords(mysqlClient, model.UpdateOption{
					Conditions:           "`ID` = ?",
					ConditionsParameters: []any{rechargeOrder.Id},
					Values: map[string]any{
						"STATUS": enum.RechargeStatus_CANCELLED.String(),
					},
				})
				continue
			}
			if confirms < 8 {
				continue
			}
			// 确认订单
			model.UpdateRechargeOrderRecords(mysqlClient, model.UpdateOption{
				Conditions:           "`ID` = ?",
				ConditionsParameters: []any{rechargeOrder.Id},
				Values: map[string]any{
					"STATUS": enum.RechargeStatus_PAID.String(),
				},
			})
			// 发起回调
			for retry := 0; retry < 5; retry++ {
				time.Sleep(time.Minute * time.Duration(retry))
				notifyModule, err := GetNotify()
				if err != nil {
					return err
				}
				err = notifyModule.Send(rechargeOrder.CallbackUrl, rechargeOrder.ExternalData)
				if err != nil {
					zap.S().Errorf("notify error: %s", err)
				} else {
					break
				}
			}
		} else {
			// 检查是否过期
			checkExpireTime(mysqlClient, rechargeOrder)
		}
	}
	return nil
}
