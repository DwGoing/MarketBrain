package config_service

import (
	context "context"

	"github.com/DwGoing/funds-system/internal/config_module"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewConfigService
type ConfigService struct {
	UnsafeConfigServiceServer

	ConfigModule *config_module.ConfigModule `singleton:""`
}

/*
@title	构造函数
@param 	service 	*ConfigService 	服务实例
@return _ 			*ConfigService 	服务实例
@return _ 			error 			异常信息
*/
func NewConfigService(service *ConfigService) (*ConfigService, error) {
	return service, nil
}

/*
@title	加载配置
@param 	Self	*ConfigService 	服务实例
@param 	ctx		context.Context 上下文
@param 	request	*LoadRequest 	请求体
@return _ 		*emptypb.Empty 	响应体
@return _ 		error 			异常信息
*/
func (Self *ConfigService) Load(ctx context.Context, request *emptypb.Empty) (*LoadResponse, error) {
	configs, err := Self.ConfigModule.Load()
	if err != nil {
		return nil, err
	}
	return &LoadResponse{
		Mnemonic:          configs.Mnemonic,
		WalletMaxNumber:   configs.WalletMaxNumber,
		ExpireTime:        configs.ExpireTime,
		ExpireDelay:       configs.ExpireDelay,
		CollectThresholds: configs.CollectThresholds,
	}, nil
}

/*
@title	修改配置
@param 	Self	*ConfigService 	服务实例
@param 	ctx		context.Context 上下文
@param 	request	*SetRequest 	请求体
@return _ 		*emptypb.Empty 	响应体
@return _ 		error 			异常信息
*/
func (Self *ConfigService) Set(ctx context.Context, request *SetRequest) (*emptypb.Empty, error) {
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
