package chain_service

import (
	context "context"
	"log"
	"os"

	"funds-system/pkg/chain_module"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewChainService
type ChainService struct {
	UnimplementedChainServiceServer

	ChainModule *chain_module.ChainModule `singleton:""`

	logger *log.Logger
}

/*
@title	构造函数
@param 	service 	*ChainService 	服务实例
@return _ 			*ChainService 	服务实例
@return _ 			error 			异常信息
*/
func NewChainService(service *ChainService) (*ChainService, error) {
	service.logger = log.New(os.Stderr, "[ChainService]", log.LstdFlags)
	return service, nil
}

/*
@title	获取某个地址的余额
@param 	Self 	*ChainService 		服务实例
@param 	ctx 	context.Context 	请求上下文
@param 	request *GetBalanceRequest 	请求体
@return _ 		*GetBalanceResponse 响应体
@return _ 		error 				异常信息
*/
func (Self *ChainService) GetBalance(ctx context.Context, request *GetBalanceRequest) (*GetBalanceResponse, error) {
	balance, err := Self.ChainModule.GetBalance(request.Address, request.Token)
	if err != nil {
		return nil, err
	}
	unconvertedBalance, err := Self.ChainModule.UnconvertValue(request.Token, balance)
	if err != nil {
		return nil, err
	}
	return &GetBalanceResponse{
		Balance: unconvertedBalance.String(),
	}, nil
}
