package bus_module

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"time"

	"funds-system/pkg/shared"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewBusModule
type BusModule struct {
	logger       *log.Logger
	RechargePaid chan shared.RechargeRecord
}

/*
@title	构造函数
@param 	module 	*BusModule 	模块实例
@return _ 		*BusModule	模块实例
@return _ 		error 		异常信息
*/
func NewBusModule(module *BusModule) (*BusModule, error) {
	module.logger = log.New(os.Stdout, "[BusModule]", log.LstdFlags)
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
	Self.logger.Printf("%s 充值完成", record.Id)
	go func() {
		retry := 0
		for {
			// 重试3次
			if retry++; retry > 3 {
				Self.logger.Printf("rechargePaidHandle Error: maximum retry limit")
				return
			}
			time.Sleep(time.Minute * time.Duration(retry))
			request, err := http.NewRequest("POST", record.CallbackUrl, bytes.NewBuffer(record.ExternalData))
			if err != nil {
				Self.logger.Printf("rechargePaidHandle Error: %s", err)
				continue
			}
			httpResponse, err := http.DefaultClient.Do(request)
			if err != nil {
				Self.logger.Printf("rechargePaidHandle Error: %s", err)
				continue
			}
			defer httpResponse.Body.Close()
			if httpResponse.StatusCode != http.StatusOK {
				Self.logger.Printf("rechargePaidHandle Error: status code not 200")
				continue
			}
			return
		}
	}()
}
