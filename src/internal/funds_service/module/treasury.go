package module

import (
	"errors"
	"strings"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/DwGoing/MarketBrain/pkg/hd_wallet"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Treasury struct{}

// @title	创建充值订单
// @param	Self				*Treasury		模块实例
// @param	externalIdentity	string			扩展标识
// @param	externalData		[]byte			扩展数据
// @param	callbackUrl			string			回调链接
// @param	chainType			enum.ChainType	主链类型
// @param	amount				float64			充值数量
// @param	walletIndex			int64			钱包索引
// @return	_					string			订单ID
// @return	_					string			接收钱包
// @return	_					error			异常信息
func (Self *Treasury) CreateRechargeOrder(
	externalIdentity string,
	externalData []byte,
	callbackUrl string,
	chainType string,
	amount float64,
	walletIndex int64,
) (string, string, time.Time, error) {
	if strings.TrimSpace(externalIdentity) == "" ||
		strings.TrimSpace(callbackUrl) == "" ||
		strings.TrimSpace(chainType) == "" ||
		amount < 1 ||
		walletIndex < 1 {
		return "", "", time.Time{}, errors.New("parameter invaild")
	}
	chain, err := new(enum.ChainType).Parse(chainType)
	if err != nil {
		return "", "", time.Time{}, err
	}
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return "", "", time.Time{}, err
	}
	chainModule, _ := GetChain()
	var walletAddress string
	switch chain {
	case enum.ChainType_TRON:
		account, err := chainModule.GetAccount(hd_wallet.Currency_TRON, walletIndex)
		if err != nil {
			return "", "", time.Time{}, err
		}
		walletAddress = account.GetAddress()
	default:
		return "", "", time.Time{}, errors.New("unsupported chain")
	}
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return "", "", time.Time{}, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return "", "", time.Time{}, err
	}
	defer db.Close()
	record := &model.RechargeOrderRecord{
		ExternalIdentity: externalIdentity,
		ExternalData:     externalData,
		CallbackUrl:      callbackUrl,
		ChainType:        chainType,
		Amount:           amount,
		WalletIndex:      walletIndex,
		WalletAddress:    walletAddress,
		ExpireAt:         time.Now().Add(time.Minute * time.Duration(config.ExpireTime)),
	}
	record, err = model.CreateRechargeOrderRecord(mysqlClient, record)
	if err != nil {
		return "", walletAddress, time.Time{}, err
	}
	return record.Id, record.WalletAddress, record.ExpireAt, nil
}

// @title	提交充值订单交易Hash
// @param	Self		*Treasury	模块实例
// @param	orderId		string		订单ID
// @param	txHash		string		交易Hash
func (Self *Treasury) SubmitRechargeOrderTransaction(orderId string, txHash string) error {
	if strings.TrimSpace(orderId) == "" ||
		strings.TrimSpace(txHash) == "" {
		return errors.New("parameter invaild")
	}
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return err
	}
	defer db.Close()
	// 检查订单
	rechargeOrders, total, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`ID` = ?",
		ConditionsParameters: []any{orderId},
	})
	if err != nil {
		return err
	}
	if total < 1 {
		return errors.New("order not existed")
	}
	// 检查TxHash
	_, total, err = model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`TX_HASH` = ?",
		ConditionsParameters: []any{txHash},
	})
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("tx hash existed")
	}
	rechargeOrder := rechargeOrders[0]
	status, _ := new(enum.RechargeStatus).Parse(rechargeOrder.Status)
	if status == enum.RechargeStatus_PAID {
		return nil
	}
	// 更新交易Hash
	err = model.UpdateRechargeOrderRecords(mysqlClient, model.UpdateOption{
		Conditions:           "`ID` = ?",
		ConditionsParameters: []any{orderId},
		Values: map[string]any{
			"TX_HASH": txHash,
			"STATUS":  enum.RechargeStatus_UNPAID.String(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// @title	取消充值订单
// @param	Self		*Treasury	模块实例
// @param	orderId		string		订单ID
// @return	_			error		异常信息
func (Self *Treasury) CancelRechargeOrder(orderId string) error {
	if strings.TrimSpace(orderId) == "" {
		return errors.New("parameter invaild")
	}
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return err
	}
	defer db.Close()
	// 更新订单状态
	err = model.UpdateRechargeOrderRecords(mysqlClient, model.UpdateOption{
		Conditions:           "`ID` = ?",
		ConditionsParameters: []any{orderId},
		Values: map[string]any{
			"STATUS": enum.RechargeStatus_CANCELLED.String(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// @title	检查充值订单状态
// @param	Self			*Treasury							模块实例
// @param	client 			*gorm.DB							mysql客户端
// @param	rechargeOrder 	*model.RechargeOrderRecord			订单
// @return	_				*model.WalletCollectionInfomation	待归集钱包
// @return	_				error								异常信息
func (Self *Treasury) checkRechargeOrderStatus(client *gorm.DB, rechargeOrder *model.RechargeOrderRecord) (*model.WalletCollectionInfomation, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return nil, err
	}
	chainType, _ := new(enum.ChainType).Parse(rechargeOrder.ChainType)
	chainConfig, ok := config.ChainConfigs[chainType.String()]
	if !ok {
		return nil, errors.New("get chain config failed")
	}
	chain, err := GetChain()
	if err != nil {
		return nil, err
	}
	// 检查交易状态
	var wallet model.WalletCollectionInfomation
	if rechargeOrder.TxHash != nil && strings.TrimSpace(*rechargeOrder.TxHash) != "" {
		tx, confirms, err := chain.DecodeTransaction(chainType, *rechargeOrder.TxHash)
		if err != nil {
			zap.S().Errorf("decode transaction error: %s", err)
			// 检查是否过期
			if rechargeOrder.ExpireAt.Before(time.Now()) {
				model.UpdateRechargeOrderRecords(client, model.UpdateOption{
					Conditions:           "`ID` = ?",
					ConditionsParameters: []any{rechargeOrder.Id},
					Values: map[string]any{
						"STATUS": enum.RechargeStatus_CANCELLED.String(),
					},
				})
				// 订单过期
				wallet.Status = enum.RechargeStatus_CANCELLED
				return &wallet, errors.New("order already expired")
			}
		}
		if !tx.Result ||
			time.UnixMilli(tx.TimeStamp).Before(rechargeOrder.CreatedAt) ||
			*tx.Contract != chainConfig.USDT ||
			tx.To != rechargeOrder.WalletAddress ||
			tx.Amount != rechargeOrder.Amount {
			model.UpdateRechargeOrderRecords(client, model.UpdateOption{
				Conditions:           "`ID` = ?",
				ConditionsParameters: []any{rechargeOrder.Id},
				Values: map[string]any{
					"STATUS": enum.RechargeStatus_CANCELLED.String(),
				},
			})
			// TxHash无效
			wallet.Status = enum.RechargeStatus_CANCELLED
			return &wallet, errors.New("tx hash invaild")
		}
		if confirms < 8 {
			wallet.Status = enum.RechargeStatus_UNPAID
			return &wallet, errors.New("insufficient number of confirmations")
		}
		// 更新订单状态
		model.UpdateRechargeOrderRecords(client, model.UpdateOption{
			Conditions:           "`ID` = ?",
			ConditionsParameters: []any{rechargeOrder.Id},
			Values: map[string]any{
				"STATUS": enum.RechargeStatus_PAID.String(),
			},
		})
		// 发起回调
		go func() {
			notifyStatus := enum.RechargeStatus_NOTIFY_OK
			for retry := 0; retry < 5; retry++ {
				time.Sleep(time.Minute * time.Duration(retry))
				notifyModule, _ := GetNotify()
				err = notifyModule.Send(rechargeOrder.CallbackUrl, rechargeOrder.ExternalData)
				if err != nil {
					notifyStatus = enum.RechargeStatus_NOTIFY_FAILED
					zap.S().Errorf("notify error: %s", err)
				} else {
					notifyStatus = enum.RechargeStatus_NOTIFY_OK
					break
				}
			}
			storageModule, _ := GetStorage()
			mysqlClient, err := storageModule.GetMysqlClient()
			if err != nil {
				return
			}
			db, err := mysqlClient.DB()
			if err != nil {
				return
			}
			defer db.Close()
			// 更新订单状态
			model.UpdateRechargeOrderRecords(mysqlClient, model.UpdateOption{
				Conditions:           "`ID` = ?",
				ConditionsParameters: []any{rechargeOrder.Id},
				Values: map[string]any{
					"STATUS": notifyStatus.String(),
				},
			})
		}()
		// 待归集
		wallet.Index = rechargeOrder.WalletIndex
		wallet.ChainType = chainType
		wallet.Address = rechargeOrder.WalletAddress
	} else {
		// 检查是否过期
		if rechargeOrder.ExpireAt.Before(time.Now()) {
			model.UpdateRechargeOrderRecords(client, model.UpdateOption{
				Conditions:           "`ID` = ?",
				ConditionsParameters: []any{rechargeOrder.Id},
				Values: map[string]any{
					"STATUS": enum.RechargeStatus_CANCELLED.String(),
				},
			})
			// 订单过期
			wallet.Status = enum.RechargeStatus_CANCELLED
			return &wallet, errors.New("order already expired")
		}
	}
	return &wallet, nil
}

// @title	检查充值订单状态
// @param	Self	*Treasury			模块实例
// @param	id		string				订单ID
// @return	_		enum.RechargeStatus	订单状态
// @return	_		error				异常信息
func (Self *Treasury) CheckRechargeOrderStatus(id string) (enum.RechargeStatus, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return 0, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	rechargeOrders, total, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`ID` = ?",
		ConditionsParameters: []any{id},
	})
	if err != nil {
		return 0, err
	}
	if total < 1 {
		return 0, errors.New("order not existed")
	}
	rechargeOrder := rechargeOrders[0]
	info, err := Self.checkRechargeOrderStatus(mysqlClient, &rechargeOrder)
	if err != nil && info == nil {
		return 0, err
	}
	return info.Status, err
}

// @title	检查充值订单状态
// @param	Self	*Treasury							模块实例
// @return	_		[]model.WalletCollectionInfomation	待归集钱包
// @return	_		error								异常信息
func (Self *Treasury) CheckRechargeOrdersStatus() ([]model.WalletCollectionInfomation, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return nil, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// 检查所有未支付订单
	rechargeOrders, _, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           "`STATUS` = ?",
		ConditionsParameters: []any{enum.RechargeStatus_UNPAID.String()},
	})
	if err != nil {
		return nil, err
	}
	wallets := []model.WalletCollectionInfomation{}
	for _, rechargeOrder := range rechargeOrders {
		wallet, err := Self.checkRechargeOrderStatus(mysqlClient, &rechargeOrder)
		if err != nil {
			zap.S().Errorf("check recharge order error: %s", err)
			continue
		}
		wallets = append(wallets, *wallet)
	}
	return wallets, nil
}

// @title	查询充值订单
// @param	Self	*Treasury					模块实例
// @return	_		[]model.RechargeOrderRecord	充值订单
// @return	_		error						异常信息
func (Self *Treasury) GetRechargeOrders(conditions string, conditionsParameters []any, pageSize int64, pageIndex int64) ([]model.RechargeOrderRecord, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return nil, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	orders, _, err := model.GetRechargeOrderRecords(mysqlClient, model.GetOption{
		Conditions:           conditions,
		ConditionsParameters: conditionsParameters,
		PageSize:             pageSize,
		PageIndex:            pageIndex,
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}
