package api

import "github.com/gin-gonic/gin"

func TreasuryRpc() *Treasury {
	api, _ := GetTreasury()
	return api
}

func TreasuryApi(router *gin.RouterGroup) {
	router.GET("ok", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
