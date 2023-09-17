package api

import (
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title	Confgig Rpc接口
func ConfigRpc() *Config {
	config, _ := GetConfig()
	return config
}

// @title	Confgig Http接口
// @param	router	*gin.RouterGroup	路由
func ConfigApi(router *gin.RouterGroup) {
	router.POST("set", SetApi)
	router.GET("load", LoadApi)
}

// @title	Treasury Rpc接口
func TreasuryRpc() *Treasury {
	treasury, _ := GetTreasury()
	return treasury
}

// @title	Treasury Http接口
// @param	router	*gin.RouterGroup	路由
func TreasuryApi(router *gin.RouterGroup) {
	router.POST("test", func(ctx *gin.Context) {
		zap.S().Warnf("Notify ===> ok")
		Response.Success(ctx, nil)
	})
	router.POST("createRechargeOrder", CreateRechargeOrderApi)
}
