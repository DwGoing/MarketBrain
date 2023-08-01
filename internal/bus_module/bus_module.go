package bus_module

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"funds-system/pkg/shared"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewBusModule
type BusModule struct {
	RechargePaid chan shared.RechargeRecord
}

/*
@title	构造函数
@param 	module 	*BusModule 	模块实例
@return _ 		*BusModule	模块实例
@return _ 		error 		异常信息
*/
func NewBusModule(module *BusModule) (*BusModule, error) {
	module.RechargePaid = make(chan shared.RechargeRecord, 1024)
	// 开启监听事件
	go module.listenEvent()
	return module, nil
}

/*
@title	监听事件
@param 	Self 	*FundsService 	服务实例
*/
func (Self *BusModule) listenEvent() {
	for {
		select {
		// 充值完成
		case record := <-Self.RechargePaid:
			Self.rechargePaidHandle(record)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

/*
@title	充值完成事件
@param 	Self 	*FundsService 			服务实例
@param 	record 	shared.RechargeRecord 	充值记录
*/
func (Self *BusModule) rechargePaidHandle(record shared.RechargeRecord) {
	log.Printf("%s 充值完成", record.Id)
	go func() {
		retry := 0
		for {
			// 重试5次
			if retry++; retry > 5 {
				log.Printf("rechargePaidHandle Error: maximum retry limit")
				return
			}
			time.Sleep(time.Minute * 2 * time.Duration(retry-1)) // 0/2/4/6/8 min
			request, err := http.NewRequest("POST", record.CallbackUrl, bytes.NewBuffer(record.ExternalData))
			if err != nil {
				log.Printf("rechargePaidHandle Error: %s", err)
				continue
			}
			httpResponse, err := http.DefaultClient.Do(request)
			if err != nil {
				log.Printf("rechargePaidHandle Error: %s", err)
				continue
			}
			defer httpResponse.Body.Close()
			if httpResponse.StatusCode != http.StatusOK {
				log.Printf("rechargePaidHandle Error: status code not 200")
				continue
			}
			return
		}
	}()
}
