package api

import (
	"github.com/DwGoing/MarketBrain/internal/funds_service/module"
	"github.com/gin-gonic/gin"
)

func TreasuryRpc() *module.Treasury {
	treasury, _ := module.GetTreasury()
	return treasury
}

func TreasuryApi(router *gin.RouterGroup) {
	treasury, _ := module.GetTreasury()
	router.POST("createRechargeOrder", treasury.CreateRechargeOrderApi)
}
