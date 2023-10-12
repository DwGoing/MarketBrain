package api

import (
	"github.com/DwGoing/MarketBrain/internal/funds_service/static/Response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RootApi(engine *gin.Engine) {
	engine.POST("test", func(ctx *gin.Context) {
		body, _ := ctx.GetRawData()
		zap.S().Warnf("Notify ===> ok %s", body)
		Response.Success(ctx, nil)
	})
	engine.POST("account", func(ctx *gin.Context) {
		Response.Success(ctx, map[string]any{
			"data": map[string]any{
				"id": 111,
			},
		})
	})
}

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
	router.POST("createRechargeOrder", CreateRechargeOrderApi)
	router.POST("submitRechargeOrderTransaction", SubmitRechargeOrderTransactionApi)
	router.POST("cancelRechargeOrder", CancelRechargeOrderApi)
	router.GET("checkRechargeOrderStatus", CheckRechargeOrderStatusApi)
}
