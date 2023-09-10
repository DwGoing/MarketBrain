package api

import (
	context "context"
	"errors"
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
=======
		return nil, err
	}
	return &CreateRechargeOrderResponse{
		OrderId:  orderId,
		ExpireAt: expireAt.String(),
	}, nil
}

type CreateRechargeOrderApiRequest struct {
	model.Request
	ExternalIdentity string  `json:"externalIdentity"`
	ExternalData     []byte  `json:"externalData"`
	CallbackUrl      string  `json:"callbackUrl"`
	ChainType        string  `json:"chainType"`
	Amount           float64 `json:"amount"`
	WalletIndex      int64   `json:"walletIndex"`
}

type CreateRechargeOrderApiResponse struct {
	model.Response
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
		ctx.JSON(200, model.Response{
			Id:      request.Id,
			Code:    enum.ApiErrorType_RequestBindError.Code(),
			Message: err.Error(),
		})
		return
	}
	if strings.TrimSpace(request.Id) == "" ||
		strings.TrimSpace(request.ExternalIdentity) == "" ||
		strings.TrimSpace(request.CallbackUrl) == "" ||
		strings.TrimSpace(request.ChainType) == "" ||
		request.Amount < 1 ||
		request.WalletIndex < 1 {
		ctx.JSON(200, model.Response{
			Id:      request.Id,
			Code:    enum.ApiErrorType_RequestBindError.Code(),
			Message: "parameter invaild",
		})
		return
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
		ctx.JSON(200, model.Response{
			Id:      request.Id,
			Code:    enum.ApiErrorType_ServiceError.Code(),
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, CreateRechargeOrderApiResponse{
		Response: model.Response{
			Id:   request.Id,
			Code: enum.ApiErrorType_Ok.Code(),
		},
		OrderId:  orderId,
		ExpireAt: expireAt,
	})
>>>>>>> 38414f3 (✨ feat: 新增创建充值订单接口)
}
