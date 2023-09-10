package module

import (
	"github.com/robfig/cron"
	"go.uber.org/zap"
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
	module.crontab.AddFunc("*/10 * * * * ?", module.checkRechargeOrderStatus)
	module.crontab.Start()
	return module, nil
}

// @title	检查充值订单状态
// @param	Self	*EventBus	模块实例
// @return	_		error		异常信息
func (Self *EventBus) checkRechargeOrderStatus() {
	treasury, _ := GetTreasury()
	err := treasury.CheckRechargeOrderStatus()
	if err != nil {
		zap.S().Errorf("check recharge order error: %s", err)
		return
	}
}
