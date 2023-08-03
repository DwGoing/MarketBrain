package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DwGoing/funds-system/internal/funds_service"
	"github.com/DwGoing/funds-system/internal/shared"

	"github.com/ahmetb/go-linq"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewFundsController
type FundsController struct {
	FundsService *funds_service.FundsService `singleton:""`
}

/*
@title	构造函数
@param 	controller 	*FundsController 	控制器实例
@return _ 			*FundsController 	控制器实例
@return _ 			error 				异常信息
*/
func NewFundsController(controller *FundsController) (*FundsController, error) {
	return controller, nil
}

type GetRechargeWalletRequest struct {
	shared.Request
	ExternalIdentity string `json:"externalIdentity,omitempty"`
	ExternalData     []byte `json:"externalData,omitempty"`
	CallbackUrl      string `json:"callbackUrl,omitempty"`
	Token            string `json:"token,omitempty"`
	Amount           string `json:"amount,omitempty"`
}

type GetRechargeWalletResponse struct {
	shared.Response
	Id       string `json:"id,omitempty"`
	Address  string `json:"address,omitempty"`
	ExpireAt int64  `json:"expireAt,omitempty"`
}

// @Summary	获取充值钱包
// @Accept	json
// @Produce	json
// @Param	request	body	GetRechargeWalletRequest	true	" "
// @Success	200  	{object}	GetRechargeWalletResponse
// @Router	/v1/funds/getRechargeWallet	[POST]
func (Self *FundsController) GetRechargeWallet(ctx *gin.Context) {
	var request GetRechargeWalletRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	response, err := Self.FundsService.GetRechargeWallet(context.Background(), &funds_service.GetRechargeWalletRequest{
		ExternalIdentity: request.ExternalIdentity,
		ExternalData:     request.ExternalData,
		CallbackUrl:      request.CallbackUrl,
		Token:            request.Token,
		Amount:           request.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, GetRechargeWalletResponse{
		Response: shared.Response{
			Code: 200,
		},
		Id:       response.Id,
		Address:  response.Address,
		ExpireAt: response.ExpireAt,
	})
}

type GetRechargeRecordsResponseResultItem struct {
	Id               string `json:"id,omitempty"`
	CreatedAt        string `json:"createdAt,omitempty"`
	UpdatedAt        string `json:"updatedAt,omitempty"`
	ExternalIdentity string `json:"externalIdentity,omitempty"`
	ExternalData     []byte `json:"externalData,omitempty"`
	CallbackUrl      string `json:"callbackUrl,omitempty"`
	Token            string `json:"token,omitempty"`
	Amount           string `json:"amount,omitempty"`
	WalletAddress    string `json:"walletAddress,omitempty"`
	Status           string `json:"status,omitempty"`
	ExpireAt         string `json:"expireAt,omitempty"`
}

type GetRechargeRecordsResponse struct {
	shared.Response
	Result []*GetRechargeRecordsResponseResultItem `json:"result,omitempty"`
	Total  int64                                   `json:"total,omitempty"`
}

// @Summary	获取充值记录
// @Produce	json
// @Param	id					query	string	false	"充值记录Id"
// @Param	externalIdentity	query	string	false	"外部标识"
// @Param	token				query	string	false	"充值Token"
// @Param	createdStart		query	int		false	"创建开始时间戳"
// @Param	createdEnd			query	int		false	"创建结束时间戳"
// @Param	pageSize			query	int		false	"页面大小"
// @Param	pageIndex			query	int		false	"页码"
// @Success	200	{object}	GetRechargeWalletResponse
// @Router	/v1/funds/getRechargeRecords 	[GET]
func (Self *FundsController) GetRechargeRecords(ctx *gin.Context) {
	// 构造查询条件
	buildConditions := func() (string, error) {
		var conditions []string
		{
			id, ok := ctx.GetQuery("id")
			if ok && len(id) > 0 {
				conditions = append(conditions, fmt.Sprintf("`ID`='%s'", id))
				goto end
			}
		}
		{
			externalIdentity, ok := ctx.GetQuery("externalIdentity")
			if ok && len(externalIdentity) > 0 {
				conditions = append(conditions, fmt.Sprintf("`EXTERNAL_IDENTITY`='%s'", externalIdentity))
				goto end
			}
		}
		{
			token, ok := ctx.GetQuery("token")
			if ok && len(token) > 0 {
				conditions = append(conditions, fmt.Sprintf("`TOKEN`='%s'", common.HexToAddress(token)))
			}
		}
		{
			createdStart, ok := ctx.GetQuery("createdStart")
			if ok && len(createdStart) > 0 {
				start, err := strconv.ParseInt(createdStart, 10, 64)
				if err != nil {
					return "", err
				}
				conditions = append(conditions, fmt.Sprintf("`CREATED_AT`>='%s'", time.Unix(start, 0)))
			}
		}
		{
			createdEnd, ok := ctx.GetQuery("createdEnd")
			if ok && len(createdEnd) > 0 {
				end, err := strconv.ParseInt(createdEnd, 10, 64)
				if err != nil {
					return "", err
				}
				conditions = append(conditions, fmt.Sprintf("`CREATED_AT`<='%s'", time.Unix(end, 0)))
			}
		}
	end:
		var where strings.Builder
		for _, item := range conditions {
			if where.Len() > 0 {
				where.WriteString(" AND ")
			}
			where.WriteString(item)
		}
		return where.String(), nil
	}
	// 分页
	var pageSize uint32
	param, ok := ctx.GetQuery("pageSize")
	if ok {
		v, err := strconv.ParseUint(param, 10, 32)
		if err != nil {
			pageSize = 20
		} else {
			pageSize = uint32(v)
		}
	}
	var pageIndex uint32
	param, ok = ctx.GetQuery("pageIndex")
	if ok {
		v, err := strconv.ParseUint(param, 10, 32)
		if err != nil {
			pageIndex = 1
		} else {
			pageIndex = uint32(v)
		}
	}
	conditions, err := buildConditions()
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	response, err := Self.FundsService.GetRechargeRecords(context.Background(), &funds_service.GetRechargeRecordsRequest{
		Conditions: conditions,
		PageSize:   pageSize,
		PageIndex:  pageIndex,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, shared.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	var result []*GetRechargeRecordsResponseResultItem
	linq.From(response.Result).SelectT(func(item *funds_service.RechargeRecord) *GetRechargeRecordsResponseResultItem {
		return &GetRechargeRecordsResponseResultItem{
			Id:               item.Id,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
			ExternalIdentity: item.ExternalIdentity,
			ExternalData:     item.ExternalData,
			CallbackUrl:      item.CallbackUrl,
			Token:            item.Token,
			Amount:           item.Amount,
			WalletAddress:    item.WalletAddress,
			Status:           item.Status,
			ExpireAt:         item.ExpireAt,
		}
	}).ToSlice(&result)
	ctx.JSON(http.StatusOK, GetRechargeRecordsResponse{
		Response: shared.Response{
			Code: 200,
		},
		Result: result,
		Total:  response.Total,
	})
}
