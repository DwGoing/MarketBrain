package controller

import (
	"context"
	"net/http"

	"github.com/DwGoing/funds-system/internal/config_service"
	"github.com/DwGoing/funds-system/pkg/shared"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewConfigController
type ConfigController struct {
	ConfigService *config_service.ConfigService `singleton:""`
}

/*
@title	构造函数
@param 	controller 	*ConfigController 	控制器实例
@return _ 			*ConfigController 	控制器实例
@return _ 			error 				异常信息
*/
func NewConfigController(controller *ConfigController) (*ConfigController, error) {
	return controller, nil
}

type LoadResponse struct {
	shared.Response
	Mnemonic          string             `json:"mnemonic,omitempty"`
	WalletMaxNumber   int64              `json:"walletMaxNumber,omitempty"`
	ExpireTime        int64              `json:"expireTime,omitempty"`
	ExpireDelay       int64              `json:"expireDelay,omitempty"`
	CollectThresholds map[string]float32 `json:"collectThresholds,omitempty"`
}

// @Summary	加载配置
// @Produce	json
// @Success	200	{object}	LoadResponse
// @Router	/v1/config/load 	[GET]
func (Self *ConfigController) Load(ctx *gin.Context) {
	response, err := Self.ConfigService.Load(context.Background(), &emptypb.Empty{})
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, &LoadResponse{
		Response: shared.Response{
			Code: 200,
		},
		Mnemonic:          response.Mnemonic,
		WalletMaxNumber:   response.WalletMaxNumber,
		ExpireTime:        response.ExpireTime,
		ExpireDelay:       response.ExpireDelay,
		CollectThresholds: response.CollectThresholds,
	})
}

type SetRequest struct {
	shared.Request
	Mnemonic          *string            `json:"mnemonic,omitempty"`
	WalletMaxNumber   *int64             `json:"walletMaxNumber,omitempty"`
	ExpireTime        *int64             `json:"expireTime,omitempty"`
	ExpireDelay       *int64             `json:"expireDelay,omitempty"`
	CollectThresholds map[string]float32 `json:"collectThresholds,omitempty"`
}

// @Summary	修改配置
// @Accept	json
// @Produce	json
// @Param	request	body	SetRequest	true	" "
// @Success	200
// @Router	/v1/config/set	[POST]
func (Self *ConfigController) Set(ctx *gin.Context) {
	var request SetRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	_, err = Self.ConfigService.Set(context.Background(), &config_service.SetRequest{
		Mnemonic:          request.Mnemonic,
		WalletMaxNumber:   request.WalletMaxNumber,
		ExpireTime:        request.ExpireTime,
		ExpireDelay:       request.ExpireDelay,
		CollectThresholds: request.CollectThresholds,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, shared.Response{
		Code: 200,
	})
}
