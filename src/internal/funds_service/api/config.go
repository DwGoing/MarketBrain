package api

import (
	context "context"

	"github.com/DwGoing/MarketBrain/internal/funds_service/api/config_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/internal/funds_service/module"
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/types/known/anypb"
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
	if request.Chains != nil {
		chains := make(map[string]any)
		for k, v := range request.Chains {
			chainType, err := new(enum.ChainType).Parse(k)
			if err != nil {
				continue
			}
			switch chainType {
			case enum.ChainType_Tron:
				chainValue, err := v.UnmarshalNew()
				if err != nil {
					continue
				}
				chain := chainValue.(*config_generated.Chain_Tron)
				chains[k] = model.Chain_Tron{
					RpcNodes:  chain.RpcNodes,
					HttpNodes: chain.HttpNodes,
					ApiKeys:   chain.ApiKeys,
				}
			default:
				continue
			}
		}
		configs["CHAINS"] = chains
	}
	if request.ExpireTime != nil {
		configs["EXPIRE_TIME"] = request.ExpireTime
	}
	if request.Mnemonic != nil {
		configs["MNEMONIC"] = request.Mnemonic
	}
	if request.WalletCollectionThreshold != nil {
		configs["WALLET_COLLECT_THRESHOLD"] = request.WalletCollectionThreshold
	}
	if request.MinGasThreshold != nil {
		configs["MIN_GAS_THRESHOLD"] = request.MinGasThreshold
	}
	if request.TransferGasAmount != nil {
		configs["TRANSFER_GAS_AMOUNT"] = request.TransferGasAmount
	}
	configModule, _ := module.GetConfig()
	err := configModule.Set(configs)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

type SetRequest struct {
	Chains                    map[string]any `json:"chains"`
	ExpireTime                *int64         `json:"expireTime"`
	Mnemonic                  *string        `json:"mnemonic"`
	WalletCollectionThreshold *float64       `json:"walletCollectionThreshold"`
	MinGasThreshold           *float64       `json:"minGasThreshold"`
	TransferGasAmount         *float64       `json:"transferGasAmount"`
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
	if request.Chains != nil {
		chains := make(map[string]any)
		for k, v := range request.Chains {
			chainType, err := new(enum.ChainType).Parse(k)
			if err != nil {
				continue
			}
			switch chainType {
			case enum.ChainType_Tron:
				var chain model.Chain_Tron
				err := mapstructure.Decode(v, &chain)
				if err != nil {
					continue
				}
				chains[k] = chain
			default:
				continue
			}
		}
		configs["CHAINS"] = chains
	}
	if request.ExpireTime != nil {
		configs["EXPIRE_TIME"] = request.ExpireTime
	}
	if request.Mnemonic != nil {
		configs["MNEMONIC"] = request.Mnemonic
	}
	if request.WalletCollectionThreshold != nil {
		configs["WALLET_COLLECT_THRESHOLD"] = request.WalletCollectionThreshold
	}
	if request.MinGasThreshold != nil {
		configs["MIN_GAS_THRESHOLD"] = request.MinGasThreshold
	}
	if request.TransferGasAmount != nil {
		configs["TRANSFER_GAS_AMOUNT"] = request.TransferGasAmount
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
	chains := make(map[string]*anypb.Any)
	for k, v := range configs.Chains {
		chainType, err := new(enum.ChainType).Parse(k)
		if err != nil {
			continue
		}
		switch chainType {
		case enum.ChainType_Tron:
			chain := v.(model.Chain_Tron)
			chains[k], _ = anypb.New(&config_generated.Chain_Tron{
				RpcNodes:  chain.RpcNodes,
				HttpNodes: chain.HttpNodes,
				ApiKeys:   chain.ApiKeys,
			})
		default:
			continue
		}
	}
	return &config_generated.LoadResponse{
		Chains:                    chains,
		ExpireTime:                configs.ExpireTime,
		Mnemonic:                  configs.Mnemonic,
		WalletCollectionThreshold: configs.WalletCollectionThreshold,
		MinGasThreshold:           configs.MinGasThreshold,
		TransferGasAmount:         configs.TransferGasAmount,
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
