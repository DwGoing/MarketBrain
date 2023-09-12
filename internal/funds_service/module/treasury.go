package module

import (
	"errors"
	"strings"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/internal/funds_service/module/treasury_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/DwGoing/MarketBrain/pkg/hd_wallet"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/DwGoing/MarketBrain/internal/funds_service/api/treasury_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/module"
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/gin-gonic/gin"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Treasury struct{}
	treasury_generated.UnimplementedTreasuryServer
}

// @title	创建充值订单
// @param	Self		*Treasury										模块实例
// @param	ctx			context.Context									上下文
// @param	request		*treasury_generated.CreateRechargeOrderRequest	请求体
// @return	_			*treasury_generated.CreateRechargeOrderResponse	响应体
// @return	_			error											异常信息
func (Self *Treasury) CreateRechargeOrderRpc(ctx context.Context, request *treasury_generated.CreateRechargeOrderRequest) (*treasury_generated.CreateRechargeOrderResponse, error) {
	treasuryModule, _ := module.GetTreasury()
	orderId, wallet, expireAt, err := treasuryModule.CreateRechargeOrder(

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
func (Self *Treasury) CreateRechargeOrder(
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
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return "", "", time.Time{}, err
	}
	chainModule, _ := GetChain()
	var walletAddress string
	switch chain {
	case enum.ChainType_TRON:
		account, err := chainModule.GetAccount(hd_wallet.Currency_TRON, walletIndex)
		if err != nil {
			return "", "", time.Time{}, err
		}
		walletAddress = account.GetAddress()
	default:
		return "", "", time.Time{}, errors.New("unsupported chain")
	}
	storageModule, _ := GetStorage()
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
		WalletIndex:      walletIndex,
		WalletAddress:    walletAddress,
		ExpireAt:         time.Now().Add(time.Minute * time.Duration(config.ExpireTime)),
	}
	record, err = model.CreateRechargeOrderRecord(mysqlClient, record)
	if err != nil {
		return "", walletAddress, time.Time{}, err
	}
	return record.Id, record.WalletAddress, record.ExpireAt, nil
}

// @title	提交充值订单交易Hash
// @param	Self		*Treasury	模块实例
// @param	orderId		string		订单ID
// @param	txHash		string		交易Hash
func (Self *Treasury) SubmitRechargeOrderTransaction(orderId string, txHash string) error {
	if strings.TrimSpace(orderId) == "" ||
		strings.TrimSpace(txHash) == "" {
		return errors.New("parameter invaild")
	}
	storageModule, _ := GetStorage()
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

// @title	取消充值订单
// @param	Self		*Treasury	模块实例
// @param	orderId		string		订单ID
// @return	_			error		异常信息
func (Self *Treasury) CancelRechargeOrder(orderId string) error {
	if strings.TrimSpace(orderId) == "" {
		return errors.New("parameter invaild")
	}
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return err
	}
	defer db.Close()
	// 更新订单状态
	err = model.UpdateRechargeOrderRecords(mysqlClient, model.UpdateOption{
		Conditions:           "`ID` = ?",
		ConditionsParameters: []any{orderId},
		Values: map[string]any{
			"STATUS": enum.RechargeStatus_CANCELLED.String(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// @title	检查充值订单状态
// @param	Self			*Treasury							模块实例
// @param	client 			*gorm.DB							mysql客户端
// @param	rechargeOrder 	*model.RechargeOrderRecord			订单
// @return	_				*model.WalletCollectionInfomation	待归集钱包
// @return	_				error								异常信息
func (Self *Treasury) checkRechargeOrderStatus(client *gorm.DB, rechargeOrder *model.RechargeOrderRecord) (*model.WalletCollectionInfomation, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return nil, err
	}
	chainType, _ := new(enum.ChainType).Parse(rechargeOrder.ChainType)
	chainConfig, ok := config.ChainConfigs[chainType.String()]
	if !ok {
		return nil, errors.New("get chain config failed")
	}
	chain, err := GetChain()
	if err != nil {
		return nil, err
	}
	// 检查交易状态
	var wallet model.WalletCollectionInfomation
	if strings.TrimSpace(rechargeOrder.TxHash) != "" {
		tx, confirms, err := chain.DecodeTransaction(chainType, rechargeOrder.TxHash)
		if err != nil {
			zap.S().Errorf("decode transaction error: %s", err)
			// 检查是否过期
			if rechargeOrder.ExpireAt.Before(time.Now()) {
				model.UpdateRechargeOrderRecords(client, model.UpdateOption{
					Conditions:           "`ID` = ?",
					ConditionsParameters: []any{rechargeOrder.Id},
					Values: map[string]any{
						"STATUS": enum.RechargeStatus_CANCELLED.String(),
					},
				})
				// 订单过期
				wallet.Status = enum.RechargeStatus_CANCELLED
				return &wallet, errors.New("order already expired")
			}
		}
		if !tx.Result ||
			tx.TimeStamp < rechargeOrder.CreatedAt.UnixMilli() ||
			tx.Contract != &chainConfig.USDT ||
			tx.From != rechargeOrder.WalletAddress ||
			tx.Amount != rechargeOrder.Amount {
				return nil, errors.New("order already expired")
			}
		}
		if !tx.Result ||
			tx.TimeStamp < rechargeOrder.CreatedAt.UnixMilli() ||
			tx.Contract != &chainConfig.USDT ||
			tx.From != rechargeOrder.WalletAddress ||
			tx.Amount != rechargeOrder.Amount {
			model.UpdateRechargeOrderRecords(client, model.UpdateOption{
				Conditions:           "`ID` = ?",
				ConditionsParameters: []any{rechargeOrder.Id},
				Values: map[string]any{
					"STATUS": enum.RechargeStatus_CANCELLED.String(),
				},
			})
			// TxHash无效
			wallet.Status = enum.RechargeStatus_CANCELLED
			return &wallet, errors.New("tx hash invaild")
		}
		if confirms < 8 {
			wallet.Status = enum.RechargeStatus_UNPAID
			return &wallet, errors.New("insufficient number of confirmations")
		}
		// 更新订单状态
		model.UpdateRechargeOrderRecords(client, model.UpdateOption{
			Conditions:           "`ID` = ?",
			ConditionsParameters: []any{rechargeOrder.Id},
			Values: map[string]any{
				"STATUS": enum.RechargeStatus_PAID.String(),
			},
		})
		// 发起回调
		notifyStatus := enum.RechargeStatus_NOTIFY_OK
		for retry := 0; retry < 5; retry++ {
			time.Sleep(time.Minute * time.Duration(retry))
			notifyModule, _ := GetNotify()
			err = notifyModule.Send(rechargeOrder.CallbackUrl, rechargeOrder.ExternalData)
			if err != nil {
				notifyStatus = enum.RechargeStatus_NOTIFY_FAILED
				zap.S().Errorf("notify error: %s", err)
			} else {
				notifyStatus = enum.RechargeStatus_NOTIFY_OK
				break
			}
		}
		// 更新订单状态
		model.UpdateRechargeOrderRecords(client, model.UpdateOption{
			Conditions:           "`ID` = ?",
			ConditionsParameters: []any{rechargeOrder.Id},
			Values: map[string]any{
				"STATUS": notifyStatus.String(),
			},
		})
		// 待归集
		wallet = model.WalletCollectionInfomation{
			Index:     rechargeOrder.WalletIndex,
			ChainType: chainType,
			Address:   rechargeOrder.WalletAddress,
			Status:    notifyStatus,
		}
	} else {
		// 检查是否过期
		if rechargeOrder.ExpireAt.Before(time.Now()) {
			model.UpdateRechargeOrderRecords(client, model.UpdateOption{
				Conditions:           "`ID` = ?",
				ConditionsParameters: []any{rechargeOrder.Id},
				Values: map[string]any{
					"STATUS": enum.RechargeStatus_CANCELLED.String(),
				},
			})
			// 订单过期
			wallet.Status = enum.RechargeStatus_CANCELLED
			return &wallet, errors.New("order already expired")
		}
	}
	return &wallet, nil
}

// @title	检查充值订单状态
// @param	Self	*Treasury			模块实例
// @param	id		string				订单ID
// @return	_		enum.RechargeStatus	订单状态
// @return	_		error				异常信息
func (Self *Treasury) CheckRechargeOrderStatus(id string) (enum.RechargeStatus, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return 0, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	rechargeOrders, total, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`ID` = ?",
		ConditionsParameters: []any{id},
	})
	if err != nil {
		return 0, err
	}
	if total < 1 {
		return 0, errors.New("order not existed")
	}
	rechargeOrder := rechargeOrders[0]
	info, err := Self.checkRechargeOrderStatus(mysqlClient, &rechargeOrder)
	if err != nil && info == nil {
		return 0, err
	}
	return info.Status, err
}

// @title	检查充值订单状态
// @param	Self	*Treasury							模块实例
// @return	_		[]model.WalletCollectionInfomation	待归集钱包
// @return	_		error								异常信息
func (Self *Treasury) CheckRechargeOrdersStatus() ([]model.WalletCollectionInfomation, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return nil, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// 检查所有未支付订单
	rechargeOrders, _, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`STATUS` = ?",
		ConditionsParameters: []any{enum.RechargeStatus_UNPAID.String()},
	})
	if err != nil {
		return nil, err
	}
	wallets := []model.WalletCollectionInfomation{}
	for _, rechargeOrder := range rechargeOrders {
		wallet, err := Self.checkRechargeOrderStatus(mysqlClient, &rechargeOrder)
		if err != nil {
			zap.S().Errorf("check recharge order error: %s", err)
			continue
		}
		wallets = append(wallets, *wallet)
	}
	return wallets, nil
}

// @title	查询充值订单
// @param	Self	*Treasury					模块实例
// @return	_		[]model.RechargeOrderRecord	充值订单
// @return	_		error						异常信息
func (Self *Treasury) GetRechargeOrders(conditions string, conditionsParameters []any, pageSize int64, pageIndex int64) ([]model.RechargeOrderRecord, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return nil, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	orders, _, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           conditions,
		ConditionsParameters: conditionsParameters,
		PageSize:             pageSize,
		PageIndex:            pageIndex,
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
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
func CreateRechargeOrderApi(ctx *gin.Context) {
	var request CreateRechargeOrderRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_RequestBindError, err)
	}
	treasuryModule, _ := module.GetTreasury()
	orderId, wallet, expireAt, err := treasuryModule.CreateRechargeOrder(
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
	}
	Response.Success(ctx, CreateRechargeOrderResponse{
		OrderId:  orderId,
		Wallet:   wallet,
		ExpireAt: expireAt,
	})
}

// @title	提交充值订单交易Hash
// @param	Self	*Treasury													模块实例
// @param	ctx		context.Context												上下文
// @param	request	*treasury_generated.SubmitRechargeOrderTransactionRequest	请求体
// @return	_		*emptypb.Empty												响应体
// @return	_		error														异常信息
func (Self *Treasury) SubmitRechargeOrderTransactionRpc(ctx context.Context, request *treasury_generated.SubmitRechargeOrderTransactionRequest) (*emptypb.Empty, error) {
	treasuryModule, _ := module.GetTreasury()
	err := treasuryModule.SubmitRechargeOrderTransaction(request.OrderId, request.TxHash)
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
func SubmitRechargeOrderTransactionApi(ctx *gin.Context) {
	var request SubmitRechargeOrderTransactionRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_RequestBindError, err)
		return
	}
	treasuryModule, _ := module.GetTreasury()
	err = treasuryModule.SubmitRechargeOrderTransaction(request.OrderId, request.TxHash)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
		return
	}
	Response.Success(ctx, nil)
}

// @title	取消充值订单
// @param	Self	*Treasury										模块实例
// @param	ctx		context.Context									上下文
// @param	request	*treasury_generated.CancelRechargeOrderRequest	请求体
// @return	_		*emptypb.Empty									响应体
// @return	_		error											异常信息
func (Self *Treasury) CancelRechargeOrderRpc(ctx context.Context, request *treasury_generated.CancelRechargeOrderRequest) (*emptypb.Empty, error) {
	treasuryModule, _ := module.GetTreasury()
	err := treasuryModule.CancelRechargeOrder(request.OrderId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

type CancelRechargeOrderRequest struct {
	OrderId string `json:"orderId"`
}

// @title	取消充值订单
// @param	Self	*Treasury										模块实例
// @param	ctx		*gin.Context									上下文
// @param	request	*treasury_generated.CancelRechargeOrderRequest	请求体
// @return	_		*emptypb.Empty									响应体
// @return	_		error											异常信息
func CancelRechargeOrderApi(ctx *gin.Context) {
	var request CancelRechargeOrderRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_RequestBindError, err)
		return
	}
	treasuryModule, _ := module.GetTreasury()
	err = treasuryModule.CancelRechargeOrder(request.OrderId)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
	}
	Response.Success(ctx, nil)
}

// @title	手动检查订单状态
// @param	Self	*Treasury												模块实例
// @param	ctx		context.Context											上下文
// @param	request	*treasury_generated.CheckRechargeOrderStatusRequest		请求体
// @return	_		*treasury_generated.CheckRechargeOrderStatusResponse	响应体
// @return	_		error													异常信息
func (Self *Treasury) CheckRechargeOrderStatusRpc(ctx context.Context, request *treasury_generated.CheckRechargeOrderStatusRequest) (*treasury_generated.CheckRechargeOrderStatusResponse, error) {
	treasuryModule, _ := module.GetTreasury()
	status, err := treasuryModule.CheckRechargeOrderStatus(request.OrderId)
	response := treasury_generated.CheckRechargeOrderStatusResponse{
		Status: treasury_generated.RechargeStatus(status),
	}
	if err != nil {
		msg := err.Error()
		response.Error = &msg
	}
	return &response, nil
}

type CheckRechargeOrderStatusResponse struct {
	Status string  `json:"orderId"`
	Error  *string `json:"error"`
}

// @title	手动检查订单状态
// @param	Self	*Treasury		模块实例
// @param	ctx		*gin.Context	上下文
// @return	_		error			异常信息
func CheckRechargeOrderStatusApi(ctx *gin.Context) {
	orderId, ok := ctx.GetQuery("orderId")
	if !ok {
		Response.Fail(ctx, enum.ApiErrorType_RequestBindError, errors.New("parameter invaild"))
		return
	}
	treasuryModule, _ := module.GetTreasury()
	status, err := treasuryModule.CheckRechargeOrderStatus(orderId)
	response := CheckRechargeOrderStatusResponse{
		Status: status.String(),
	}
	if err != nil {
		msg := err.Error()
		response.Error = &msg
	}
	Response.Success(ctx, response)
		OrderId:  orderId,
		ExpireAt: expireAt.String(),
	}, nil
}

type CreateRechargeOrderApiRequest struct {
	ExternalIdentity string  `json:"externalIdentity"`
	ExternalData     []byte  `json:"externalData"`
	CallbackUrl      string  `json:"callbackUrl"`
	ChainType        string  `json:"chainType"`
	Amount           float64 `json:"amount"`
	WalletIndex      int64   `json:"walletIndex"`
}

type CreateRechargeOrderApiResponse struct {
	OrderId  string    `json:"orderId"`
	ExpireAt time.Time `json:"expireAt"`
}

// @title	创建充值订单
// @param	Self	*Treasury		服务实例
// @param	ctx		*gin.Context	上下文
func (Self *Treasury) CreateRechargeOrderApi(ctx *gin.Context) {
	var request CreateRechargeOrderApiRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_RequestBindError, err)
	}
	if strings.TrimSpace(request.ExternalIdentity) == "" ||
		strings.TrimSpace(request.CallbackUrl) == "" ||
		strings.TrimSpace(request.ChainType) == "" ||
		request.Amount < 1 ||
		request.WalletIndex < 1 {
		Response.Fail(ctx, enum.ApiErrorType_ParameterError, err)
	}
	orderId, expireAt, err := Self.createRechargeOrder(
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
	Response.Success(ctx, CreateRechargeOrderApiResponse{
		OrderId:  orderId,
		ExpireAt: expireAt,
	})
}

// @title	检查充值订单状态
// @param	Self	*Treasury	服务实例
// @return	_		error		异常信息
func (Self *Treasury) CheckRechargeOrderStatus() error {
	zap.S().Errorf("check recharge order error: %s", "xxxxxxx")
	return nil
	// client, err := Self.Storage.GetMysqlClient()
	// if err != nil {
	// 	zap.S().Errorf("check recharge order error: %s", err)
	// 	return
	// }
}
