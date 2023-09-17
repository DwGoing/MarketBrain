package api

import (
	context "context"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/api/treasury_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/module"
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/gin-gonic/gin"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Treasury struct {
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
