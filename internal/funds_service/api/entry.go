package api

import (
	"github.com/DwGoing/MarketBrain/internal/funds_service/module"
	"github.com/gin-gonic/gin"
)

// @title	Confgig Rpc接口
func ConfigRpc() *module.Config {
	config, _ := module.GetConfig()
	return config
}

// @title	Confgig Http接口
// @param	router	*gin.RouterGroup	路由
func ConfigApi(router *gin.RouterGroup) {
	config, _ := module.GetConfig()
	router.POST("set", config.SetApi)
	router.GET("load", config.LoadApi)
}

// @title	Treasury Rpc接口
func TreasuryRpc() *module.Treasury {
	treasury, _ := module.GetTreasury()
	return treasury
}

// @title	Treasury Http接口
// @param	router	*gin.RouterGroup	路由
func TreasuryApi(router *gin.RouterGroup) {
	treasury, _ := module.GetTreasury()
	router.POST("createRechargeOrder", treasury.CreateRechargeOrderApi)
}
