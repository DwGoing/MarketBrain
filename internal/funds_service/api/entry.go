package api

import "github.com/gin-gonic/gin"

func TreasuryRpc() *Treasury {
	api, _ := GetTreasury()
	return api
}

func TreasuryApi(router *gin.RouterGroup) {
	api, _ := GetTreasury()
	router.POST("createRechargeOrder", api.CreateRechargeOrderApi)
}
