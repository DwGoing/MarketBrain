package data_provider

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

type BinanceDataProvider struct {
	websocketConn *websocket.Conn
	sendLock      *sync.RWMutex
	subscriptions map[SubscribeType]BinanceSubscription
	requests      map[string]SubscribeType
}

func (Self SubscribeType) Method() (string, error) {
	switch Self {
	case SubscribeType_Kline:
		return "klines", nil
	default:
		return "", errors.New("unknow type")
	}
}

type BinanceSubscription struct {
	Request  BinanceWebsocketRequest
	Crontab  *cron.Cron
	Callback func(any)
}

type BinanceWebsocketRequest struct {
	Id       string       `json:"id"`
	Method   string       `json:"method"`
	Params   any          `json:"params,omitempty"`
	Callback func([]byte) `json:"-"`
}

type BinanceWebsocketError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type BinanceWebsocketRateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
	Count         int    `json:"count"`
}

type BinanceWebsocketResponse struct {
	Id         string                      `json:"id"`
	Status     int                         `json:"status"`
	Error      BinanceWebsocketError       `json:"error"`
	RateLimits []BinanceWebsocketRateLimit `json:"rateLimits"`
}

type BinanceWebsocketKlineResponse struct {
	BinanceWebsocketResponse
	Result [][]any `json:"result"`
}

func NewBinanceDataProvider() (IDataProvider, error) {
	var dataProvider = BinanceDataProvider{}
	err := dataProvider.startWebsocket()
	if err != nil {
		return nil, err
	}
	dataProvider.Subscribe(SubscribeType_Ping, nil, "*/3 * * * * ?", nil)
	return &dataProvider, nil
}

func (Self *BinanceDataProvider) startWebsocket() error {
	// 连接websocket
	for {
		zap.S().Debug("==============> 连接Websocket")
		conn, _, err := websocket.DefaultDialer.Dial("wss://ws-api.binance.com:443/ws-api/v3", nil)
		if err != nil {
			zap.S().Errorf("connect error: %s", err)
			time.Sleep(time.Second * 5)
			continue
		}
		Self.websocketConn = conn
		break
	}
	// 接收消息
	go func() {
		defer Self.websocketConn.Close()
		for {
			_, msgBytes, err := Self.websocketConn.ReadMessage()
			if err != nil {
				zap.S().Errorf("receive message error: %s", err)
				go Self.startWebsocket()
				return
			}
			zap.S().Debugf("receive message: %s", msgBytes)
			var response BinanceWebsocketResponse
			err = json.Unmarshal(msgBytes, &response)
			if err != nil {
				zap.S().Debugf("can not unmarshal message: %s", msgBytes)
				continue
			}
			subscribeType, ok := Self.requests[response.Id]
			if ok {
				switch subscribeType {
				case SubscribeType_Ping:
					log.Println("ping")
				case SubscribeType_Kline:
					var response BinanceWebsocketKlineResponse
					err = json.Unmarshal(msgBytes, &response)
					if err != nil {
						zap.S().Debugf("can not unmarshal message: %s", msgBytes)
						continue
					}
				default:
					zap.S().Debugf("unhandle message: %s", msgBytes)
				}
				delete(Self.requests, response.Id)
			} else {
				zap.S().Debugf("unhandle message: %s", msgBytes)
			}
		}
	}()
	return nil
}

func (Self *BinanceDataProvider) websocketSendMessage(request BinanceWebsocketRequest) error {
	request.Id = uuid.NewString()
	requestBytes, _ := json.Marshal(request)
	Self.sendLock.Lock()
	defer Self.sendLock.Unlock()
	if Self.websocketConn == nil {
		return nil
	}
	err := Self.websocketConn.WriteMessage(1, requestBytes)
	if err != nil {
		return err
	}
	return nil
}

func (Self *BinanceDataProvider) Subscribe(subscribeType SubscribeType, parameters any, interval string, callback func(any)) error {
	request := BinanceWebsocketRequest{
		Id: uuid.NewString(),
	}
	method, err := subscribeType.Method()
	if err != nil {
		return nil
	}
	request.Method = method
	subscription := BinanceSubscription{
		Request:  request,
		Crontab:  cron.New(),
		Callback: callback,
	}
	subscription.Crontab.AddFunc(interval, func() {
		Self.websocketSendMessage(request)
		Self.requests[request.Id] = subscribeType
	})
	subscription.Crontab.Start()
	Self.subscriptions[subscribeType] = subscription
	return nil
}

func (Self *BinanceDataProvider) Unsubscribe(subscribeType SubscribeType) error {
	subscription, ok := Self.subscriptions[subscribeType]
	if ok {
		subscription.Crontab.Stop()
	}
	delete(Self.subscriptions, subscribeType)
	return nil
}
