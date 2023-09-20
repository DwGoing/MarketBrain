package module

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/DwGoing/MarketBrain/pkg/hd_wallet"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Chain struct{}

// @title	获取子钱包
// @param	Self			*Chain				模块实例
// @param	currencyType 	hd_wallet.Currency	币种类型
// @param	index			int64				钱包索引
// @return	_				*hd_wallet.Account	子钱包
// @return	_				error				异常信息
func (Self *Chain) GetAccount(currencyType hd_wallet.Currency, index int64) (*hd_wallet.Account, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return nil, err
	}
	hdWallet, err := hd_wallet.FromMnemonic(config.Mnemonic, "")
	if err != nil {
		return nil, err
	}
	account, err := hdWallet.GetAccount(currencyType, index)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// @title	获取当前高度
// @param	Self		*Chain			模块实例
// @param	chainType	enum.ChainType	链类型
// @return	_			int64			交易信息
// @return	_			error			异常信息
func (Self *Chain) GetCurrentHeight(chainType enum.ChainType) (int64, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return 0, err
	}
	chainConfig, ok := config.ChainConfigs[chainType.String()]
	if !ok || len(chainConfig.Nodes) < 1 {
		return 0, errors.New("no chain config")
	}
	var height int64
	switch chainType {
	case enum.ChainType_TRON:
		tron, _ := GetTron()
		client, err := tron.GetTronClient(chainConfig.Nodes, chainConfig.ApiKey)
		if err != nil {
			return 0, err
		}
		block, err := client.GetNowBlock()
		if err != nil {
			return 0, err
		}
		height = block.BlockHeader.RawData.Number
	default:
		return 0, errors.New("unsupported chain type")
	}
	return height, nil
}

// @title	获取钱包余额
// @param	Self		*Chain			模块实例
// @param	chainType	enum.ChainType	链类型
// @param	contract	*string			合约地址
// @param	wallet		string			钱包地址
// @return	_			float64			余额
// @return	_			error			异常信息
func (Self *Chain) GetBalance(chainType enum.ChainType, contract *string, wallet string) (float64, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return 0, err
	}
	chainConfig, ok := config.ChainConfigs[chainType.String()]
	if !ok || len(chainConfig.Nodes) < 1 {
		return 0, errors.New("no chain config")
	}
	var balance float64
	switch chainType {
	case enum.ChainType_TRON:
		tron, _ := GetTron()
		client, err := tron.GetTronClient(chainConfig.Nodes, chainConfig.ApiKey)
		if err != nil {
			return 0, err
		}
		if contract == nil {
			account, err := client.GetAccount(wallet)
			if err != nil {
				return 0, err
			}
			balance, _ = new(big.Float).Quo(new(big.Float).SetInt64(account.Balance), big.NewFloat(1e6)).Float64()
		} else {
			balanceBigInt, err := client.TRC20ContractBalance(wallet, *contract)
			if err != nil {
				return 0, err
			}
			decimalsBigInt, err := client.TRC20GetDecimals(*contract)
			if err != nil {
				return 0, err
			}
			balance, _ = new(big.Float).Quo(new(big.Float).SetInt(balanceBigInt), big.NewFloat(math.Pow10(int(decimalsBigInt.Int64())))).Float64()
		}
	default:
		return 0, errors.New("unsupported chain type")
	}
	return balance, nil
}

// @title	发送代币
// @param	Self		*Chain				模块实例
// @param	chainType	enum.ChainType		链类型
// @param	token		*string				合约地址
// @param	from		*hd_wallet.Account	发送账户
// @param	to			string				接收地址
// @param	amount		float64				数量
// @param	remarks		string				备注
// @return	_			error				异常信息
func (Self *Chain) Transfer(chainType enum.ChainType, token *string, from *hd_wallet.Account, to string, amount float64, remarks string) (string, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return "", err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	txHash, err := func(chainType enum.ChainType, token *string, from *hd_wallet.Account, to string, amount float64) (string, error) {
		configModule, _ := GetConfig()
		config, err := configModule.Load()
		if err != nil {
			return "", err
		}
		chainConfig, ok := config.ChainConfigs[chainType.String()]
		if !ok || len(chainConfig.Nodes) < 1 {
			return "", errors.New("no chain config")
		}
		var txHash string
		switch chainType {
		case enum.ChainType_TRON:
			tron, _ := GetTron()
			client, err := tron.GetTronClient(chainConfig.Nodes, chainConfig.ApiKey)
			if err != nil {
				return "", err
			}
			var tx *api.TransactionExtention
			if token == nil {
				amountInt64, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e6)).Int64()
				tx, err = client.Transfer(from.GetAddress(), to, amountInt64)
				if err != nil {
					return "", err
				}
			} else {
				decimalsBigInt, err := client.TRC20GetDecimals(*token)
				if err != nil {
					return "", err
				}
				amountBigInt, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(math.Pow10(int(decimalsBigInt.Int64())))).Int(new(big.Int))
				tx, err = client.TRC20Send(from.GetAddress(), to, *token, amountBigInt, 300000000)
				if err != nil {
					return "", err
				}
			}
			txInfo, err := tron.SendTronTransaction(client, from.PrivateKey.ToECDSA(), tx.Transaction, true)
			if err != nil {
				return "", err
			}
			txHash = common.Bytes2Hex(txInfo.GetId())
		default:
			return "", errors.New("unsupported chain type")
		}
		return txHash, nil
	}(chainType, token, from, to, amount)
	// 创建日志
	record := model.TransferRecord{
		ChainType:   chainType.String(),
		Token:       token,
		FromIndex:   from.Index,
		FromAddress: from.GetAddress(),
		To:          to,
		Amount:      amount,
		Status:      enum.TransferStatus_SUCCESS.String(),
		Remarks:     remarks,
	}
	if err != nil {
		record.Status = enum.TransferStatus_FAILED.String()
		record.Error = err.Error()
	}
	_, err = model.CreateTransferRecord(mysqlClient, &record)
	if err != nil {
		return "", err
	}
	return txHash, err
}

// @title	查询块中交易
// @param	Self		*Chain				模块实例
// @param	chainType	enum.ChainType		链类型
// @param	start		int64				数量
// @param	end			int64				备注
// @return	_			[]model.Transaction	异常信息
// @return	_			error				异常信息
func (Self *Chain) GetTransactionFromBlocks(chainType enum.ChainType, start int64, end int64) ([]model.Transaction, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return nil, err
	}
	chainConfig, ok := config.ChainConfigs[chainType.String()]
	if !ok || len(chainConfig.Nodes) < 1 {
		return nil, errors.New("no chain config")
	}
	var result []model.Transaction
	switch chainType {
	case enum.ChainType_TRON:
		tron, _ := GetTron()
		client, err := tron.GetTronClient(chainConfig.Nodes, chainConfig.ApiKey)
		if err != nil {
			return nil, err
		}
		result, err = tron.GetTransactionsFromBlocks(client, start, end)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported chain type")
	}
	return result, nil
}

// @title	解析交易
// @param	Self		*Chain				模块实例
// @param	chainType	enum.ChainType		链类型
// @param	txHash		string				交易Hash
// @return	_			*model.Transaction	交易信息
// @return	_			error				异常信息
func (Self *Chain) DecodeTransaction(chainType enum.ChainType, txHash string) (*model.Transaction, int64, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return nil, 0, err
	}
	chainConfig, ok := config.ChainConfigs[chainType.String()]
	if !ok || len(chainConfig.Nodes) < 1 {
		return nil, 0, errors.New("no chain config")
	}
	var transaction model.Transaction
	var confirms int64
	switch chainType {
	case enum.ChainType_TRON:
		tron, _ := GetTron()
		client, err := tron.GetTronClient(chainConfig.Nodes, chainConfig.ApiKey)
		if err != nil {
			return nil, 0, err
		}
		// 最新区块
		block, err := client.GetNowBlock()
		if err != nil {
			return nil, 0, err
		}
		tx, err := tron.DecodeTransaction(client, txHash)
		if err != nil {
			return nil, 0, err
		}
		transaction = *tx
		confirms = block.BlockHeader.RawData.Number - transaction.Height
	default:
		if err != nil {
			return nil, 0, errors.New("unsupported chain type")
		}
	}
	return &transaction, confirms, err
}

// @title	获取Tron客户端
// @param	Self		*Chain			模块实例
// @param	config	*ChainConfig		链配置
// @return	_		*client.GrpcClient	客户端
// @return	_		error				异常信息
func (Self *Chain) getTronClient(config model.ChainConfig) (*client.GrpcClient, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(config.Nodes))))
	if err != nil {
		return nil, err
	}
	grpcClient := client.NewGrpcClient(config.Nodes[index.Int64()])
	err = grpcClient.SetAPIKey(config.ApiKey)
	if err != nil {
		return nil, err
	}
	err = grpcClient.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return grpcClient, nil
}

// @title	发送Tron交易
// @param	Self		*Chain					模块实例
// @param	client		*client.GrpcClient		客户端
// @param	privateKey	*ecdsa.PrivateKey		私钥
// @param	tx			*core.Transaction		交易
// @param	waitReceipt	*client.GrpcClient		是否等待结果
// @return	_			*core.TransactionInfo	交易信息
// @return	_			error					异常信息
func (Self *Chain) sendTronTransaction(client *client.GrpcClient, privateKey *ecdsa.PrivateKey, tx *core.Transaction, waitReceipt bool) (*core.TransactionInfo, error) {
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}
	tx.Signature = append(tx.Signature, signature)
	result, err := client.Broadcast(tx)
	if err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("bad transaction: %v", string(result.GetMessage()))
	}
	var transaction *core.TransactionInfo
	start := 0
	for {
		if start++; start > 10 {
			return nil, errors.New("transaction info not found")
		}
		transaction, err = client.GetTransactionInfoByID(common.BytesToHexString(hash))
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if transaction.Result != 0 {
			return nil, errors.New(string(transaction.ResMessage))
		}
		break
	}
	return transaction, err
}

// @title	解析交易
// @param	Self		*Chain			模块实例
// @param	chainType	enum.ChainType	链类型
// @param	contract	*string			合约地址
// @param	txHash		string			交易Hash
// @return	_			bool			交易状态
// @return	_			int64			时间戳
// @return	_			string			收款地址
// @return	_			float64			金额
// @return	_			int64			确认数
// @return	_			error			异常信息
func (Self *Chain) DecodeTransaction(chainType enum.ChainType, contract *string, txHash string) (bool, int64, string, float64, int64, error) {
	var (
		result    bool
		timeStamp int64
		to        string
		amount    float64
		confirms  int64
		err       error
	)
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return result, timeStamp, to, amount, confirms, err
	}
	chainConfig, ok := config.ChainConfigs[chainType.String()]
	if !ok || len(chainConfig.Nodes) < 1 {
		return result, timeStamp, to, amount, confirms, errors.New("no chain config")
	}
	switch chainType {
	case enum.ChainType_TRON:
		client, err := Self.getTronClient(chainConfig)
		if err != nil {
			return result, timeStamp, to, amount, confirms, err
		}
		tx, err := client.GetTransactionInfoByID(txHash)
		if err != nil {
			return result, timeStamp, to, amount, confirms, err
		}
		receiptResult := tx.GetReceipt().GetResult()
		result = receiptResult == core.Transaction_Result_DEFAULT || receiptResult == core.Transaction_Result_SUCCESS
		if !result {
			return result, timeStamp, to, amount, confirms, err
		}
		timeStamp = tx.GetBlockTimeStamp()
		if contract == nil {
			txWithRawData, err := client.GetTransactionByID(txHash)
			if err != nil {
				return result, timeStamp, to, amount, confirms, err
			}
			var tc core.TransferContract
			err = txWithRawData.RawData.GetContract()[0].GetParameter().UnmarshalTo(&tc)
			if err != nil {
				return result, timeStamp, to, amount, confirms, err
			}
			to = common.EncodeCheck(tc.GetToAddress())
			amount, _ = new(big.Float).Quo(new(big.Float).SetInt64(tc.GetAmount()), big.NewFloat(1e6)).Float64()
		} else {
			if *contract != common.EncodeCheck(tx.GetContractAddress()) {
				return result, timeStamp, to, amount, confirms, errors.New("contract not match")
			}
			log := tx.GetLog()[0]
			if common.BytesToHexString(log.GetTopics()[0]) != common.BytesToHexString(common.Keccak256([]byte("Transfer(address,address,uint256)"))) {
				return result, timeStamp, to, amount, confirms, errors.New("function not match")
			}
			to = common.EncodeCheck(append([]byte{0x41}, log.GetTopics()[2][12:]...))
			decimalsBigInt, err := client.TRC20GetDecimals(*contract)
			if err != nil {
				return result, timeStamp, to, amount, confirms, err
			}
			amount, _ = new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).SetBytes(log.Data)), big.NewFloat(math.Pow10(int(decimalsBigInt.Int64())))).Float64()
		}
		lastestBlock, err := client.GetNowBlock()
		if err != nil {
			return result, timeStamp, to, amount, confirms, err
		}
		confirms = lastestBlock.BlockHeader.RawData.Number - tx.BlockNumber
	default:
		err = errors.New("unsupported chain type")
	}
	return result, timeStamp, to, amount, confirms, err
}

// @title	获取钱包余额
// @param	Self		*Chain			模块实例
// @param	chainType	enum.ChainType	链类型
// @param	contract	*string			合约地址
// @param	wallet		string			钱包地址
// @return	_			float64			余额
// @return	_			error			异常信息
func (Self *Chain) GetBalance(chainType enum.ChainType, contract *string, wallet string) (float64, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return 0, err
	}
	chainConfig, ok := config.ChainConfigs[chainType.String()]
	if !ok || len(chainConfig.Nodes) < 1 {
		return 0, errors.New("no chain config")
	}
	var balance float64
	switch chainType {
	case enum.ChainType_TRON:
		client, err := Self.getTronClient(chainConfig)
		if err != nil {
			return 0, err
		}
		if contract == nil {
			account, err := client.GetAccount(wallet)
			if err != nil {
				return 0, err
			}
			balance, _ = new(big.Float).Quo(new(big.Float).SetInt64(account.Balance), big.NewFloat(1e6)).Float64()
		} else {
			balanceBigInt, err := client.TRC20ContractBalance(wallet, *contract)
			if err != nil {
				return 0, err
			}
			decimalsBigInt, err := client.TRC20GetDecimals(*contract)
			if err != nil {
				return 0, err
			}
			balance, _ = new(big.Float).Quo(new(big.Float).SetInt(balanceBigInt), big.NewFloat(math.Pow10(int(decimalsBigInt.Int64())))).Float64()
		}
	default:
		return 0, errors.New("unsupported chain type")
	}
	return balance, nil
}

// @title	发送代币
// @param	Self		*Chain				模块实例
// @param	chainType	enum.ChainType		链类型
// @param	token		*string				合约地址
// @param	from		*hd_wallet.Account	发送账户
// @param	to			string				接收地址
// @param	amount		float64				数量
// @param	remarks		string				备注
// @return	_			error				异常信息
func (Self *Chain) Transfer(chainType enum.ChainType, token *string, from *hd_wallet.Account, to string, amount float64, remarks string) (string, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return "", err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	txHash, err := func(chainType enum.ChainType, token *string, from *hd_wallet.Account, to string, amount float64) (string, error) {
		configModule, _ := GetConfig()
		config, err := configModule.Load()
		if err != nil {
			return "", err
		}
		chainConfig, ok := config.ChainConfigs[chainType.String()]
		if !ok || len(chainConfig.Nodes) < 1 {
			return "", errors.New("no chain config")
		}
		var txHash string
		switch chainType {
		case enum.ChainType_TRON:
			client, err := Self.getTronClient(chainConfig)
			if err != nil {
				return "", err
			}
			var tx *api.TransactionExtention
			if token == nil {
				amountInt64, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e6)).Int64()
				tx, err = client.Transfer(from.GetAddress(), to, amountInt64)
				if err != nil {
					return "", err
				}
			} else {
				decimalsBigInt, err := client.TRC20GetDecimals(*token)
				if err != nil {
					return "", err
				}
				amountBigInt, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(math.Pow10(int(decimalsBigInt.Int64())))).Int(new(big.Int))
				tx, err = client.TRC20Send(from.GetAddress(), to, *token, amountBigInt, 300000000)
				if err != nil {
					return "", err
				}
			}
			txInfo, err := Self.sendTronTransaction(client, from.PrivateKey.ToECDSA(), tx.Transaction, true)
			if err != nil {
				return "", err
			}
			txHash = common.Bytes2Hex(txInfo.GetId())
		default:
			return "", errors.New("unsupported chain type")
		}
		return txHash, nil
	}(chainType, token, from, to, amount)
	// 创建日志
	record := model.TransferRecord{
		ChainType:   chainType.String(),
		Token:       token,
		FromIndex:   from.Index,
		FromAddress: from.GetAddress(),
		To:          to,
		Amount:      amount,
		Status:      enum.TransferStatus_SUCCESS.String(),
		Remarks:     remarks,
	}
	if err != nil {
		record.Status = enum.TransferStatus_FAILED.String()
		record.Error = err.Error()
	}
	_, err = model.CreateTransferRecord(mysqlClient, &record)
	if err != nil {
		return "", err
	}
	return txHash, err
}
