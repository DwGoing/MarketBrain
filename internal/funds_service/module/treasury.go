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
	"github.com/gin-gonic/gin"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Treasury struct {
	treasury_generated.UnimplementedTreasuryServer

	Storage *Storage `normal:""`
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
// @return	_					error			异常信息
func (Self *Treasury) createRechargeOrder(
	externalIdentity string,
	externalData []byte,
	callbackUrl string,
	chainType string,
	amount float64,
	walletIndex int64,
) (string, time.Time, error) {
	_, err := new(enum.ChainType).Parse(chainType)
	if err != nil {
		return "", time.Time{}, err
	}
	client, err := Self.Storage.GetMysqlClient()
	if err != nil {
		return "", time.Time{}, err
	}
	record := &model.RechargeOrderRecord{
		ExternalIdentity: externalIdentity,
		ExternalData:     externalData,
		CallbackUrl:      callbackUrl,
		ChainType:        chainType,
		Amount:           amount,
		WalletIndex:      walletIndex,
		WalletAddress:    "0x",
		ExpireAt:         time.Now(),
	}
	record, err = model.CreateRechargeOrderRecord(client, record)
	if err != nil {
		return "", time.Time{}, err
	}
	return record.Id, record.ExpireAt, nil
}

// @title	创建充值订单
// @param	Self		*Treasury										服务实例
// @param	ctx			context.Context									上下文
// @param	request		*treasury_generated.CreateRechargeOrderRequest	请求体
// @return	_			*treasury_generated.CreateRechargeOrderResponse	响应体
// @return	_			error											异常信息
func (Self *Treasury) CreateRechargeOrder(ctx context.Context, request *treasury_generated.CreateRechargeOrderRequest) (*treasury_generated.CreateRechargeOrderResponse, error) {
	if strings.TrimSpace(request.ExternalIdentity) == "" ||
		strings.TrimSpace(request.CallbackUrl) == "" ||
		strings.TrimSpace(request.ChainType) == "" ||
		request.Amount < 1 ||
		request.WalletIndex < 1 {
		return nil, errors.New("parameter invaild")
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
		return nil, err
	}
	return &treasury_generated.CreateRechargeOrderResponse{
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
