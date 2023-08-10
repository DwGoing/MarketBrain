package funds_service

import (
	"context"
	"net/http"

	"github.com/DwGoing/OnlyPay/internal/shared"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
func LoadConfig(ctx *gin.Context) {
	fundsService, err := GetFundsServiceSingleton()
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	response, err := fundsService.LoadConfig(context.Background(), &emptypb.Empty{})
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
func SetConfig(ctx *gin.Context) {
	var request SetConfigRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	fundsService, err := GetFundsServiceSingleton()
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	_, err = fundsService.SetConfig(context.Background(), &SetConfigRequest{
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
