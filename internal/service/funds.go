package service

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/api"
	"github.com/DwGoing/MarketBrain/internal/funds_service/api/config_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/api/treasury_generated"
	"github.com/DwGoing/MarketBrain/internal/funds_service/module"
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/alibaba/ioc-golang/extension/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewFunds
type Funds struct {
	RpcPort  *config.ConfigInt `config:",service.rpc"`
	HttpPort *config.ConfigInt `config:",service.http"`
	EventBus *module.EventBus  `singleton:""`
}

// @title	构造函数
// @param 	service *FundsService 	服务实例
// @return _ 		*FundsService 	服务实例
// @return _ 		error 			异常信息
func NewFunds(service *Funds) (*Funds, error) {
	// 初始化Rpc
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", service.RpcPort.Value()))
		if err != nil {
			zap.S().Fatalf("监听器初始化失败：%s", err)
		}
		server := grpc.NewServer()
		config_generated.RegisterConfigServer(server, api.ConfigRpc())
		treasury_generated.RegisterTreasuryServer(server, api.TreasuryRpc())
		zap.S().Infof("RPC服务正在监听: %v", listener.Addr())
		if err = server.Serve(listener); err != nil {
			zap.S().Errorf("RPC服务开启失败: %s", err)
		}
	}()
	// 初始化Http
	go func() {
		engine := gin.Default()
		// 验证RequestId
		engine.Use(func(ctx *gin.Context) {
			requestId, ok := ctx.GetQuery("requestId")
			if !ok || strings.TrimSpace(requestId) == "" {
				Response.Fail(ctx, enum.ApiErrorType_ParameterError, errors.New("request id invaild"))
			}
			ctx.Set("requestId", requestId)
		})
		config := engine.Group("config")
		api.ConfigApi(config)
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
