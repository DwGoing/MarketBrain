package api

import (
	context "context"

	"github.com/DwGoing/MarketBrain/internal/funds_service/api/config_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/internal/funds_service/module"
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/gin-gonic/gin"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Config struct {
	config_generated.UnimplementedConfigServer
}

// @title	更新配置
// @param	Self		*Config							模块实例
// @param	ctx			context.Context					上下文
// @param	request		*config_generated.SetRequest	请求体
// @return	_			*emptypb.Empty					响应体
// @return	_			error							异常信息
func (Self *Config) SetRpc(ctx context.Context, request *config_generated.SetRequest) (*emptypb.Empty, error) {
	configs := make(map[string]any)
	if request.Mnemonic != nil {
		configs["MNEMONIC"] = request.Mnemonic
	}
	if request.ChainConfigs != nil {
		for k := range request.ChainConfigs {
			_, err := new(enum.ChainType).Parse(k)
			if err != nil {
				continue
			}
		}
		configs["CHAIN_CONFIGS"] = request.ChainConfigs
	}
	configModule, _ := module.GetConfig()
	err := configModule.Set(configs)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

type SetRequest struct {
	Mnemonic     *string                      `json:"mnemonic"`
	ChainConfigs map[string]model.ChainConfig `json:"chainConfigs"`
}

// @title	更新配置
// @param	Self	*Config			模块实例
// @param	ctx		*gin.Context	上下文
func SetApi(ctx *gin.Context) {
	var request SetRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_RequestBindError, err)
	}
	configs := make(map[string]any)
	if request.Mnemonic != nil {
		configs["MNEMONIC"] = request.Mnemonic
	}
	if request.ChainConfigs != nil {
		for k := range request.ChainConfigs {
			_, err := new(enum.ChainType).Parse(k)
			if err != nil {
				continue
			}
		}
		configs["CHAIN_CONFIGS"] = request.ChainConfigs
	}
	configModule, _ := module.GetConfig()
	err = configModule.Set(configs)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
	}
	Response.Success(ctx, nil)
}

// @title	加载配置
// @param	Self		*Config							模块实例
// @param	ctx			context.Context					上下文
// @param	request		*mptypb.Empty					请求体
// @return	_			*config_generated.LoadResponse	响应体
// @return	_			error							异常信息
func (Self *Config) LoadRpc(ctx context.Context, request *emptypb.Empty) (*config_generated.LoadResponse, error) {
	configModule, _ := module.GetConfig()
	configs, err := configModule.Load()
	if err != nil {
		return nil, err
	}
	chainConfigs := make(map[string]*config_generated.ChainConfig)
	for k, v := range configs.ChainConfigs {
		chainConfigs[k] = &config_generated.ChainConfig{
			Nodes: v.Nodes,
		}
	}
	return &config_generated.LoadResponse{
		Mnemonic:     configs.Mnemonic,
		ChainConfigs: chainConfigs,
	}, nil
}

// @title	加载配置
// @param	Self	*Config			模块实例
// @param	ctx		*gin.Context	上下文
func LoadApi(ctx *gin.Context) {
	configModule, _ := module.GetConfig()
	configs, err := configModule.Load()
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
	}
	Response.Success(ctx, configs)
}
