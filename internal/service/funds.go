package service

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/api"
	"github.com/alibaba/ioc-golang/extension/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewFunds
type Funds struct {
	RpcPort  *config.ConfigInt `config:",service.rpc"`
	HttpPort *config.ConfigInt `config:",service.http"`

	crontab     *cron.Cron
	redisClient *redis.Client
	mysqlClient *gorm.DB
}

// @title	构造函数
// @param 	service *FundsService 	服务实例
// @return _ 		*FundsService 	服务实例
// @return _ 		error 			异常信息
func NewFunds(service *Funds) (*Funds, error) {
	service.crontab = cron.New()
	service.crontab.Start()
	// 初始化Rpc
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", service.RpcPort.Value()))
		if err != nil {
			zap.S().Fatalf("监听器初始化失败：%s", err)
		}
		server := grpc.NewServer()
		api.RegisterTreasuryServer(server, api.TreasuryRpc())
		zap.S().Infof("RPC服务正在监听: %v", listener.Addr())
		if err = server.Serve(listener); err != nil {
			zap.S().Errorf("RPC服务开启失败: %s", err)
		}
	}()
	// 初始化Http
	go func() {
		engine := gin.Default()
		treasury := engine.Group("treasury")
		api.TreasuryApi(treasury)
		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", service.HttpPort.Value()),
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		zap.S().Infof("HTTP服务正在监听: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			zap.S().Errorf("HTTP服务开启失败: %s", err)
		}
	}()
	return service, nil
}
