package user_service

import (
	"log"
	"os"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewUserService
type UserService struct {
	logger *log.Logger
}

/*
@title	构造函数
@param 	service *UserService 	服务实例
@return _ 		*UserService 	服务实例
@return _ 		error 			异常信息
*/
func NewUserService(service *UserService) (*UserService, error) {
	service.logger = log.New(os.Stdout, "[UserService]", log.LstdFlags)
	return service, nil
}

func (Self *UserService) Initialize() error {
	return nil
}
