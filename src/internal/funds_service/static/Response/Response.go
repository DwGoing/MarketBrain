package Response

import (
	"net/http"

	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Id      string `json:"id"`
	Code    int64  `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{
		Id:      ctx.GetString("requestId"),
		Code:    enum.ApiErrorType_Ok.Code(),
		Message: "success",
		Data:    data,
	})
	ctx.Abort()
}

func Fail(ctx *gin.Context, errorType enum.ApiErrorType, err error) {
	ctx.JSON(http.StatusOK, Response{
		Id:      ctx.GetString("requestId"),
		Code:    errorType.Code(),
		Message: errorType.String(err),
	})
	ctx.Abort()
}
