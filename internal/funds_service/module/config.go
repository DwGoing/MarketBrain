package module

import (
	context "context"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/internal/funds_service/module/config_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Config struct {
	config_generated.UnimplementedConfigServer

	Storage *Storage `normal:""`
}

type Configs struct {
	Mnemonic     string                 `mapstructure:"MNEMONIC" json:"mnemonic"`
	ChainConfigs map[string]ChainConfig `mapstructure:"CHAIN_CONFIGS" json:"chainConfigs"`
}

type ChainConfig struct {
	USDT   string   `json:"usdt"`
	Nodes  []string `json:"nodes"`
	ApiKey string   `json:"apiKey"`
}

// @title	更新配置
// @param	Self	*Config		模块实例
// @return	_		*Configs	配置
// @return	_		error		异常信息
func (Self *Config) set(configs map[string]any) error {
	configRecords := []model.ConfigRecord{}
	for k, v := range configs {
		configRecords = append(configRecords, model.ConfigRecord{
			Key:   k,
			Value: v,
		})
	}
	mysqlClient, err := Self.Storage.GetMysqlClient()
	if err != nil {
		return err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = model.UpdateConfigRecords(mysqlClient, configRecords)
	if err != nil {
		return err
	}
	return nil
}

// @title	更新配置
// @param	Self		*Config							服务实例
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
	err := Self.set(configs)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

type SetRequest struct {
	Mnemonic     *string                `json:"mnemonic"`
	ChainConfigs map[string]ChainConfig `json:"chainConfigs"`
}

// @title	更新配置
// @param	Self	*Config			服务实例
// @param	ctx		*gin.Context	上下文
func (Self *Config) SetApi(ctx *gin.Context) {
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
	err = Self.set(configs)
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
	}
	Response.Success(ctx, nil)
}

// @title	加载配置
// @param	Self	*Config		模块实例
// @return	_		*Configs	配置
// @return	_		error		异常信息
func (Self *Config) load() (*Configs, error) {
	mysqlClient, err := Self.Storage.GetMysqlClient()
	if err != nil {
		return nil, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	configRecords, err := model.GetConfigRecords(mysqlClient)
	if err != nil {
		return nil, err
	}
	configMap := map[string]any{}
	linq.From(configRecords).ToMapByT(&configMap, func(item model.ConfigRecord) string {
		return item.Key
	}, func(item model.ConfigRecord) any {
		return item.Value
	})
	var configs *Configs
	mapstructure.Decode(configMap, &configs)
	return configs, nil
}

// @title	加载配置
// @param	Self	*Config		模块实例
// @return	_		*Configs	配置
// @return	_		error		异常信息
func (Self *Config) Load() (*Configs, error) {
	return Self.load()
}

// @title	加载配置
// @param	Self		*Config							服务实例
// @param	ctx			context.Context					上下文
// @param	request		*mptypb.Empty					请求体
// @return	_			*config_generated.LoadResponse	响应体
// @return	_			error							异常信息
func (Self *Config) LoadRpc(ctx context.Context, request *emptypb.Empty) (*config_generated.LoadResponse, error) {
	configs, err := Self.load()
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
// @param	Self	*Config			服务实例
// @param	ctx		*gin.Context	上下文
func (Self *Config) LoadApi(ctx *gin.Context) {
	configs, err := Self.load()
	if err != nil {
		Response.Fail(ctx, enum.ApiErrorType_ServiceError, err)
	}
	Response.Success(ctx, configs)
}
