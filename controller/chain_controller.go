package controller

import (
	"context"
	"net/http"

	"github.com/DwGoing/funds-system/internal/chain_service"
	"github.com/DwGoing/funds-system/internal/shared"

	"github.com/gin-gonic/gin"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewChainController
type ChainController struct {
	ChainService *chain_service.ChainService `singleton:""`
}

/*
@title	构造函数
@param 	controller 	*ChainController 	控制器实例
@return _ 			*ChainController 	控制器实例
@return _ 			error 				异常信息
*/
func NewChainController(controller *ChainController) (*ChainController, error) {
	return controller, nil
}

type GetBalanceRequest struct {
	shared.Request
	Address string `json:"address,omitempty"`
}

type GetBalanceResponse struct {
	shared.Response
	Balance string `json:"balance,omitempty"`
}

// @Summary	查询余额
// @Produce	json
// @Param	address	query	string	true	"查询地址"
// @Param	token	query	string	false	"查询Token"
// @Success	200	{object}	GetBalanceResponse
// @Router	/v1/chain/getBalance 	[GET]
func (Self *ChainController) GetBalance(ctx *gin.Context) {
	address, ok := ctx.GetQuery("address")
	if !ok || len(address) <= 0 {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: "address invaild",
		})
		return
	}
	token, _ := ctx.GetQuery("token")
	response, err := Self.ChainService.GetBalance(context.Background(), &chain_service.GetBalanceRequest{
		Address: address,
		Token:   token,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, &GetBalanceResponse{
		Response: shared.Response{
			Code: 200,
		},
		Balance: response.Balance,
	})
}
