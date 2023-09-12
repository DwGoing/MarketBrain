package module

import (
	"github.com/robfig/cron"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewEventBus
type EventBus struct {
	crontab *cron.Cron
}

// @title	构造函数
// @param 	module *Storage 	模块实例
// @return _ 		*Storage 	模块实例
// @return _ 		error 		异常信息
func NewEventBus(module *EventBus) (*EventBus, error) {
	module.crontab = cron.New()
	module.crontab.Start()
	return module, nil
}
