package service

import (
	"github.com/alibaba/ioc-golang/extension/config"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewFundsService
type FundsService struct {
	X *config.ConfigInt `config:",xxx"`
}

func NewFundsService(service *FundsService) (*FundsService, error) {
	return service, nil
}
