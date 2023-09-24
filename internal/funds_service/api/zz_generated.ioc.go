//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package api

import (
	contextx "context"
	"github.com/DwGoing/MarketBrain/internal/funds_service/api/config_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/api/treasury_generated"
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	util "github.com/alibaba/ioc-golang/autowire/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &config_{}
		},
	})
	configStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &Config{}
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
	}
	normal.RegisterStructDescriptor(configStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &treasury_{}
		},
	})
	treasuryStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &Treasury{}
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
	}
	normal.RegisterStructDescriptor(treasuryStructDescriptor)
}

type config_ struct {
	SetRpc_  func(ctx contextx.Context, request *config_generated.SetRequest) (*emptypb.Empty, error)
	LoadRpc_ func(ctx contextx.Context, request *emptypb.Empty) (*config_generated.LoadResponse, error)
}

func (c *config_) SetRpc(ctx contextx.Context, request *config_generated.SetRequest) (*emptypb.Empty, error) {
	return c.SetRpc_(ctx, request)
}

func (c *config_) LoadRpc(ctx contextx.Context, request *emptypb.Empty) (*config_generated.LoadResponse, error) {
	return c.LoadRpc_(ctx, request)
}

type treasury_ struct {
	CreateRechargeOrderRpc_            func(ctx contextx.Context, request *treasury_generated.CreateRechargeOrderRequest) (*treasury_generated.CreateRechargeOrderResponse, error)
	SubmitRechargeOrderTransactionRpc_ func(ctx contextx.Context, request *treasury_generated.SubmitRechargeOrderTransactionRequest) (*emptypb.Empty, error)
	CancelRechargeOrderRpc_            func(ctx contextx.Context, request *treasury_generated.CancelRechargeOrderRequest) (*emptypb.Empty, error)
	CancelRechargeOrderApi_            func(ctx *gin.Context)
	CheckRechargeOrderStatusRpc_       func(ctx contextx.Context, request *treasury_generated.CheckRechargeOrderStatusRequest) (*treasury_generated.CheckRechargeOrderStatusResponse, error)
}

func (t *treasury_) CreateRechargeOrderRpc(ctx contextx.Context, request *treasury_generated.CreateRechargeOrderRequest) (*treasury_generated.CreateRechargeOrderResponse, error) {
	return t.CreateRechargeOrderRpc_(ctx, request)
}

func (t *treasury_) SubmitRechargeOrderTransactionRpc(ctx contextx.Context, request *treasury_generated.SubmitRechargeOrderTransactionRequest) (*emptypb.Empty, error) {
	return t.SubmitRechargeOrderTransactionRpc_(ctx, request)
}

func (t *treasury_) CancelRechargeOrderRpc(ctx contextx.Context, request *treasury_generated.CancelRechargeOrderRequest) (*emptypb.Empty, error) {
	return t.CancelRechargeOrderRpc_(ctx, request)
}

func (t *treasury_) CancelRechargeOrderApi(ctx *gin.Context) {
	t.CancelRechargeOrderApi_(ctx)
}

func (t *treasury_) CheckRechargeOrderStatusRpc(ctx contextx.Context, request *treasury_generated.CheckRechargeOrderStatusRequest) (*treasury_generated.CheckRechargeOrderStatusResponse, error) {
	return t.CheckRechargeOrderStatusRpc_(ctx, request)
}

type ConfigIOCInterface interface {
	SetRpc(ctx contextx.Context, request *config_generated.SetRequest) (*emptypb.Empty, error)
	LoadRpc(ctx contextx.Context, request *emptypb.Empty) (*config_generated.LoadResponse, error)
}

type TreasuryIOCInterface interface {
	CreateRechargeOrderRpc(ctx contextx.Context, request *treasury_generated.CreateRechargeOrderRequest) (*treasury_generated.CreateRechargeOrderResponse, error)
	SubmitRechargeOrderTransactionRpc(ctx contextx.Context, request *treasury_generated.SubmitRechargeOrderTransactionRequest) (*emptypb.Empty, error)
	CancelRechargeOrderRpc(ctx contextx.Context, request *treasury_generated.CancelRechargeOrderRequest) (*emptypb.Empty, error)
	CancelRechargeOrderApi(ctx *gin.Context)
	CheckRechargeOrderStatusRpc(ctx contextx.Context, request *treasury_generated.CheckRechargeOrderStatusRequest) (*treasury_generated.CheckRechargeOrderStatusResponse, error)
}

var _configSDID string

func GetConfig() (*Config, error) {
	if _configSDID == "" {
		_configSDID = util.GetSDIDByStructPtr(new(Config))
	}
	i, err := normal.GetImpl(_configSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*Config)
	return impl, nil
}

func GetConfigIOCInterface() (ConfigIOCInterface, error) {
	if _configSDID == "" {
		_configSDID = util.GetSDIDByStructPtr(new(Config))
	}
	i, err := normal.GetImplWithProxy(_configSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(ConfigIOCInterface)
	return impl, nil
}

var _treasurySDID string

func GetTreasury() (*Treasury, error) {
	if _treasurySDID == "" {
		_treasurySDID = util.GetSDIDByStructPtr(new(Treasury))
	}
	i, err := normal.GetImpl(_treasurySDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*Treasury)
	return impl, nil
}

func GetTreasuryIOCInterface() (TreasuryIOCInterface, error) {
	if _treasurySDID == "" {
		_treasurySDID = util.GetSDIDByStructPtr(new(Treasury))
	}
	i, err := normal.GetImplWithProxy(_treasurySDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(TreasuryIOCInterface)
	return impl, nil
}
