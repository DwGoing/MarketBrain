package funds_service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DwGoing/OnlyPay/internal/shared"
	"github.com/ahmetb/go-linq"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// @Summary	获取充值钱包
// @Accept	json
// @Produce	json
// @Param	request body	GetRechargeWalletRequest	true	" "
// @Success	200  	{object}	GetRechargeWalletResponse
// @Router	/v1/funds/getRechargeWallet	[POST]
func GetRechargeWallet(ctx *gin.Context) {
	var request GetRechargeWalletRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
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
	response, err := fundsService.GetRechargeWallet(context.Background(), &request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":     200,
		"id":       response.Id,
		"address":  response.Address,
		"expireAt": response.ExpireAt,
	})
}

type GetRechargeRecordsResultItem struct {
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
func GetRechargeRecords(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
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
	response, err := fundsService.GetRechargeRecords(context.Background(), &GetRechargeRecordsRequest{
		Conditions: conditions,
		PageSize:   pageSize,
		PageIndex:  pageIndex,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	var result []GetRechargeRecordsResultItem
	linq.From(response.Result).SelectT(func(item *RechargeRecord) GetRechargeRecordsResultItem {
		return GetRechargeRecordsResultItem{
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
	ctx.JSON(http.StatusOK, gin.H{
		"code":   200,
		"result": result,
		"total":  response.Total,
	})
}
