package module

import (
	"context"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/DwGoing/MarketBrain/pkg/hd_wallet"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewEventBus
type EventBus struct {
	crontab              *cron.Cron
	walletCollectChannel chan model.WalletCollectionInfomation

	currentHeights map[enum.ChainType]int64
}

// @title	构造函数
// @param 	module	*EventBus 	模块实例
// @return _ 		*EventBus 	模块实例
// @return _ 		error 		异常信息
func NewEventBus(module *EventBus) (*EventBus, error) {
	chainModule, err := GetChain()
	if err != nil {
		return nil, err
	}
	// 初始化区块高度
	module.currentHeights = map[enum.ChainType]int64{
		enum.ChainType_TRON: 0,
	}
	for k := range module.currentHeights {
		height, err := chainModule.GetCurrentHeight(k)
		if err != nil {
			return nil, err
		}
		module.currentHeights[k] = height
	}

	module.crontab = cron.New(cron.WithSeconds(), cron.WithChain(cron.DelayIfStillRunning(cron.DefaultLogger)))
	_, err = module.crontab.AddFunc("*/10 * * * * ?", module.checkRechargeOrderStatus)
	if err != nil {
		return nil, err
	}
	_, err = module.crontab.AddFunc("*/10 * * * * ?", module.checkNewTransaction)
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
				// 验证钱包是否为当前子钱包
				var currency hd_wallet.Currency
				var chainConfig model.ChainConfig
				switch wallet.ChainType {
				case enum.ChainType_TRON:
					config, ok := config.ChainConfigs[wallet.ChainType.String()]
					if !ok {
						zap.S().Errorf("no chain config")
						return
					}
					currency = hd_wallet.Currency_TRON
					chainConfig = config
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
				mainAccount, err := chainModule.GetAccount(currency, 0)
				if err != nil {
					zap.S().Errorf("get account error: %s", err)
					return
				}
				// 补充Gas
				if gasBalance < config.MinGasThreshold {
					_, err = chainModule.Transfer(enum.ChainType_TRON, nil, mainAccount, wallet.Address, config.TransferGasAmount, "补充Gas")
					if err != nil {
						zap.S().Errorf("transfer gas error: %s", err)
						return
					}
				}
				// 检查余额
				balance, err := chainModule.GetBalance(wallet.ChainType, &chainConfig.USDT, wallet.Address)
				if err != nil {
					zap.S().Errorf("get balance error: %s", err)
					return
				}
				// 钱包归集
				if balance >= config.WalletCollectionThreshold {
					txHash, err := chainModule.Transfer(enum.ChainType_TRON, &chainConfig.USDT, account, mainAccount.GetAddress(), balance, "钱包归集")
					if err != nil {
						zap.S().Errorf("transfer usdt error: %s", err)
						return
					}
					zap.S().Debugf("collect [%s] %s === %f ===> %s", txHash, wallet.Address, balance, mainAccount.GetAddress())
				}
			}(configModule, chainModule)
		default:
			time.Sleep(time.Second * 5)
		}
	}
}

// @title	交易监听
// @param	Self	*EventBus	模块实例
func (Self *EventBus) checkNewTransaction() {
	storageModule, _ := GetStorage()
	redisClient, err := storageModule.GetRedisClient()
	if err != nil {
		zap.S().Errorf("get redis client error: %s", err)
		return
	}
	defer redisClient.Close()
	// 加锁
	lock := "NEW_TRANSACTION_CHECKING"
	ok, err := redisClient.SetNX(context.Background(), lock, "", time.Duration(time.Minute*10)).Result()
	if err != nil {
		zap.S().Errorf("get NEW_TRANSACTION_CHECKING lock error: %s", err)
		return
	}
	if !ok {
		zap.S().Errorf("get NEW_TRANSACTION_CHECKING lock failed")
		return
	}
	// 解锁
	defer redisClient.Del(context.Background(), lock).Result()
	zap.S().Debugf("start check new transaction")
	chainModule, _ := GetChain()
	treasury, _ := GetTreasury()
	for k, v := range Self.currentHeights {
		// 查询当前高度
		height, err := chainModule.GetCurrentHeight(k)
		if err != nil {
			continue
		}
		start := v
		var end int64
		// 单次最多查询5个区块
		if height-v > 5 {
			end = height + 5
		} else {
			end = height
		}
		// 获取区块中代币交易
		transactions, err := chainModule.GetTransactionFromBlocks(k, start, end)
		if err != nil {
			continue
		}
		// 查找匹配的订单
		err = func() error {
			for _, transaction := range transactions {
				orders, err := treasury.GetRechargeOrders(
					"`CREATED_AT` <= ? AND `CHAIN_TYPE` = ? AND `AMOUNT` = ? AND `WALLET_ADDRESS` = ? AND `STATUS` = ?",
					[]any{time.UnixMilli(transaction.TimeStamp), transaction.ChainType.String(), transaction.Amount, transaction.To, enum.RechargeStatus_UNPAID.String()},
					1, 1,
				)
				if err != nil {
					return err
				}
				if len(orders) < 1 {
					continue
				}
				order := orders[0]
				// 提交Hash等待验证
				err = treasury.SubmitRechargeOrderTransaction(order.Id, transaction.Hash)
				if err != nil {
					return err
				}
			}
			return nil
		}()
		if err != nil {
			zap.S().Errorf("order match error: %s", err)
			continue
		}
		Self.currentHeights[k] = end
	}
}
