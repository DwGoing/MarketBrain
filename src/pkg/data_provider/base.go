package data_provider

type SubscribeType uint8

const (
	SubscribeType_Ping  SubscribeType = 1
	SubscribeType_Kline SubscribeType = 2
)

type IDataProvider interface {
	Subscribe(subscribeType SubscribeType, parameters any, interval string, callback func(any)) error
	Unsubscribe(subscribeType SubscribeType) error
}

type KlineResponse struct {
}
