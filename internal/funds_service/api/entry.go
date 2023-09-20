package api

import (
	"github.com/gin-gonic/gin"
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
	router.POST("createRechargeOrder", CreateRechargeOrderApi)
}
