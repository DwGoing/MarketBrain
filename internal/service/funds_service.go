package service

import (
	"github.com/alibaba/ioc-golang/extension/config"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewFundsService
type FundsService struct {
	RpcPort               *config.ConfigInt    `config:",service.rpc"`
	HttpPort              *config.ConfigInt    `config:",service.http"`
	RedisConnectionString *config.ConfigString `config:",storage.redis"`
	MysqlConnectionString *config.ConfigString `config:",storage.mysql"`
}

func NewFundsService(service *FundsService) (*FundsService, error) {
	return service, nil
}
