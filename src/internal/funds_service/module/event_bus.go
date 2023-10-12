package module

import (
	"context"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/DwGoing/MarketBrain/pkg/hd_wallet"
	"github.com/ahmetb/go-linq"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewEventBus
type EventBus struct {
	crontab              *cron.Cron
	walletCollectChannel chan model.WalletCollectionInfomation
}

// @title	构造函数
// @param 	module	*EventBus 	模块实例
// @return _ 		*EventBus 	模块实例
// @return _ 		error 		异常信息
func NewEventBus(module *EventBus) (*EventBus, error) {
	module.crontab = cron.New(cron.WithSeconds(), cron.WithChain(cron.DelayIfStillRunning(cron.DefaultLogger)))
	_, err := module.crontab.AddFunc("*/10 * * * * ?", module.checkRechargeOrderStatus)
	if err != nil {
		return nil, err
	}
	_, err = module.crontab.AddFunc("*/10 * * * * ?", module.checkTronNewTransaction)
	if err != nil {
		return nil, err
	}
	module.crontab.Start()
	module.walletCollectChannel = make(chan model.WalletCollectionInfomation, 1024)
	go module.collectWallet()
	return module, nil
}

// @title	检查充值订单状态
// @param	Self	*EventBus	模块实例
func (Self *EventBus) checkRechargeOrderStatus() {
	storageModule, _ := GetStorage()
	redisClient, err := storageModule.GetRedisClient()
	if err != nil {
		zap.S().Errorf("get redis client error: %s", err)
		return
	}
	defer redisClient.Close()
	// 加锁
	lock := "RECHARGE_ORDER_STATUS_CHEAKING"
	ok, err := redisClient.SetNX(context.Background(), lock, "", time.Duration(time.Minute*10)).Result()
	if err != nil {
		zap.S().Errorf("get RECHARGE_ORDER_STATUS_CHEAKING lock error: %s", err)
		return
	}
	if !ok {
		zap.S().Errorf("get RECHARGE_ORDER_STATUS_CHEAKING lock failed")
		return
	}
	// 解锁
	defer redisClient.Del(context.Background(), lock).Result()
	zap.S().Debugf("start check recharge order status")
	treasury, _ := GetTreasury()
	wallets, err := treasury.CheckRechargeOrdersStatus()
	if err != nil {
		zap.S().Errorf("check recharge order error: %s", err)
		return
	}
	// 待归集
	for _, wallet := range wallets {
		Self.walletCollectChannel <- wallet
	}
}

// @title	钱包归集
// @param	Self	*EventBus			模块实例
// @param	wallets	map[int64]string	待检查的钱包
func (Self *EventBus) collectWallet() {
	defer func() {
		if rec := recover(); rec != nil {
			zap.S().Errorf("panic error: %s", rec)
			Self.collectWallet()
		}
	}()
	for {
		configModule, _ := GetConfig()
		chainModule, _ := GetChain()
		select {
		case wallet := <-Self.walletCollectChannel:
			go func(configModule *Config, chainModule *Chain) {
				config, err := configModule.Load()
				if err != nil {
					zap.S().Errorf("get config error: %s", err)
					return
				}
				_ = config
				// 验证钱包是否为当前子钱包
				var currency hd_wallet.Currency
				// var chainConfig model.Chain
				switch wallet.ChainType {
				case enum.ChainType_Tron:
					// config, ok := config.ChainConfigs[wallet.ChainType.String()]
					// if !ok {
					// 	zap.S().Errorf("no chain config")
					// 	return
					// }
					// currency = hd_wallet.Currency_TRON
					// chainConfig = config
				default:
					return
				}
				account, err := chainModule.GetAccount(currency, wallet.Index)
				if err != nil {
					zap.S().Errorf("get account error: %s", err)
					return
				}
				expectAddress := account.GetAddress()
				if wallet.Address != expectAddress {
					zap.S().Errorf("index[%d] and address[%s] not match, expect: %s", wallet.Index, wallet.Address, expectAddress)
					return
				}
				// 检查Gas余额
				gasBalance, err := chainModule.GetBalance(wallet.ChainType, nil, wallet.Address)
				if err != nil && err.Error() != "account not found" {
					zap.S().Errorf("get gas balance error: %s", err)
					return
				}
				_ = gasBalance
				mainAccount, err := chainModule.GetAccount(currency, 0)
				if err != nil {
					zap.S().Errorf("get account error: %s", err)
					return
				}
				_ = mainAccount
				// 补充Gas
				// if gasBalance < config.MinGasThreshold {
				// 	_, err = chainModule.Transfer(enum.ChainType_Tron, nil, mainAccount, wallet.Address, config.TransferGasAmount, "补充Gas")
				// 	if err != nil {
				// 		zap.S().Errorf("transfer gas error: %s", err)
				// 		return
				// 	}
				// }
				// 检查余额
				// balance, err := chainModule.GetBalance(wallet.ChainType, &chainConfig.USDT, wallet.Address)
				// if err != nil {
				// 	zap.S().Errorf("get balance error: %s", err)
				// 	return
				// }
				// 钱包归集
				// if balance >= config.WalletCollectionThreshold {
				// 	txHash, err := chainModule.Transfer(enum.ChainType_Tron, &chainConfig.USDT, account, chainConfig.CollectionTarget, balance, "钱包归集")
				// 	if err != nil {
				// 		zap.S().Errorf("transfer usdt error: %s", err)
				// 		return
				// 	}
				// 	zap.S().Debugf("collect [%s] %s === %f ===> %s", txHash, wallet.Address, balance, chainConfig.CollectionTarget)
				// }
			}(configModule, chainModule)
		default:
			time.Sleep(time.Second * 5)
		}
	}
}

// @title	交易监听
// @param	Self	*EventBus	模块实例
func (Self *EventBus) checkTronNewTransaction() {
	storageModule, _ := GetStorage()
	redisClient, err := storageModule.GetRedisClient()
	if err != nil {
		zap.S().Errorf("get redis client error: %s", err)
		return
	}
	defer redisClient.Close()
	// 加锁
	lock := "NEW_TRON_TRANSACTION_CHECKING"
	ok, err := redisClient.SetNX(context.Background(), lock, "", time.Duration(time.Minute*10)).Result()
	if err != nil {
		zap.S().Errorf("get %s lock error: %s", lock, err)
		return
	}
	if !ok {
		zap.S().Errorf("get %s lock failed", lock)
		return
	}
	// 解锁
	defer redisClient.Del(context.Background(), lock).Result()
	treasury, _ := GetTreasury()
	// 查询所有UNPAID的订单
	orders, err := treasury.GetRechargeOrders(
		"`STATUS` = ?",
		[]any{enum.RechargeStatus_UNPAID.String()},
		10000, 1,
	)
	if err != nil {
		zap.S().Errorf("get UNPAID order error: %s", err)
		return
	}
	if len(orders) < 1 {
		return
	}
	// 根据接收地址分组
	groups := make(map[string][]model.RechargeOrderRecord)
	linq.From(orders).GroupByT(
		func(item model.RechargeOrderRecord) string {
			return item.WalletAddress
		}, func(item model.RechargeOrderRecord) model.RechargeOrderRecord {
			return item
		},
	).ToMapByT(&groups, func(item linq.Group) string {
		return item.Key.(string)
	}, func(item linq.Group) []model.RechargeOrderRecord {
		var values []model.RechargeOrderRecord
		linq.From(item.Group).SelectT(func(item interface{}) model.RechargeOrderRecord {
			return item.(model.RechargeOrderRecord)
		}).ToSlice(&values)
		return values
	})
	// configModule, _ := GetConfig()
	// config, err := configModule.Load()
	// if err != nil {
	// 	zap.S().Errorf("get config error: %s", err)
	// 	return
	// }
	// chainConfig, ok := config.ChainConfigs[enum.ChainType_Tron.String()]
	// if !ok {
	// 	zap.S().Errorf("no chain config")
	// 	return
	// }
	for address, groupOrders := range groups {
		go func(address string, orders []model.RechargeOrderRecord) {
			// 查询近期交易
			earliestTime := linq.From(orders).SelectT(func(item model.RechargeOrderRecord) time.Time {
				return item.CreatedAt
			}).OrderByT(func(item time.Time) int64 {
				return item.UnixMilli()
			}).First().(time.Time)
			chainModule, _ := GetChain()

			var transactions []model.Transaction
			_ = earliestTime
			_ = chainModule

			// transactions, err := chainModule.GetTransactionsByAddress(enum.ChainType_Tron, address, &chainConfig.USDT, earliestTime)
			if err != nil {
				zap.S().Errorf("get transactions error: %s", err)
				return
			}
			for _, transaction := range transactions {
				// 查找匹配的订单
				matchedOrders, err := treasury.GetRechargeOrders(
					"`CREATED_AT` <= ? AND `CHAIN_TYPE` = ? AND `AMOUNT` = ? AND `WALLET_ADDRESS` = ? AND (`TX_HASH` IS NULL OR `TX_HASH` = '')",
					[]any{time.UnixMilli(transaction.TimeStamp), transaction.ChainType.String(), transaction.Amount, transaction.To},
					100, 1,
				)
				if err != nil {
					zap.S().Errorf("get order error: %s", err)
					continue
				}
				if len(matchedOrders) < 1 {
					continue
				}
				// 选择最早的订单
				order := linq.From(matchedOrders).OrderByDescendingT(func(item model.RechargeOrderRecord) int64 {
					return item.CreatedAt.UnixMilli()
				}).First().(model.RechargeOrderRecord)
				// 提交Hash等待验证
				err = treasury.SubmitRechargeOrderTransaction(order.Id, transaction.Hash)
				if err != nil {
					zap.S().Errorf("submit order transaction error: %s", err)
					continue
				}
				zap.S().Infof("found matched order: %s ===> %s", order.Id, transaction.Hash)
			}
		}(address, groupOrders)
	}
}
