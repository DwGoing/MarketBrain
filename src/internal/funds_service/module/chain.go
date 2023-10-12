package module

import (
	"errors"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/DwGoing/MarketBrain/pkg/hd_wallet"
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
	var height int64
	switch chainType {
	case enum.ChainType_Tron:
		tron, _ := GetTron()
		h, err := tron.GetCurrentHeight()
		if err != nil {
			return 0, err
		}
		height = h
	default:
		return 0, errors.New("unsupported chain type")
	}
	return height, nil
}

// @title	获取钱包余额
// @param	Self		*Chain			模块实例
// @param	chainType	enum.ChainType	链类型
// @param	token	*string			合约地址
// @param	address		string			钱包地址
// @return	_			float64			余额
// @return	_			error			异常信息
func (Self *Chain) GetBalance(chainType enum.ChainType, token *string, address string) (float64, error) {
	var balance float64
	switch chainType {
	case enum.ChainType_Tron:
		tron, _ := GetTron()
		b, err := tron.GetBalance(token, address)
		if err != nil {
			return 0, err
		}
		balance = b
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
		var txHash string
		switch chainType {
		case enum.ChainType_Tron:
			tron, _ := GetTron()
			hash, err := tron.Transfer(token, from, to, amount)
			if err != nil {
				return "", err
			}
			txHash = hash
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

// @title	查询块信息
// @param	Self		*Chain				模块实例
// @param	chainType	enum.ChainType		链类型
// @param	number		int64				块高
// @return	_			*model.Block		块信息
// @return	_			error				异常信息
func (Self *Chain) GetBlock(chainType enum.ChainType, number int64) (*model.Block, error) {
	var result model.Block
	switch chainType {
	case enum.ChainType_Tron:
		tron, _ := GetTron()
		block, err := tron.GetBlockByNumber(number)
		if err != nil {
			return nil, err
		}
		result = *block
	default:
		return nil, errors.New("unsupported chain type")
	}
	return &result, nil
}

// @title	查询块中交易
// @param	Self		*Chain				模块实例
// @param	chainType	enum.ChainType		链类型
// @param	start		int64				数量
// @param	end			int64				备注
// @return	_			[]model.Transaction	异常信息
// @return	_			error				异常信息
func (Self *Chain) GetTransactionFromBlocks(chainType enum.ChainType, start int64, end int64) ([]model.Transaction, error) {
	var result []model.Transaction
	switch chainType {
	case enum.ChainType_Tron:
		tron, _ := GetTron()
		transactions, err := tron.GetTronTransactionsFromBlocks(start, end)
		if err != nil {
			return nil, err
		}
		result = transactions
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
	var transaction model.Transaction
	var confirms int64
	switch chainType {
	case enum.ChainType_Tron:
		tron, _ := GetTron()
		client, err := tron.GetTronRpcClient()
		if err != nil {
			return nil, 0, err
		}
		// 最新区块
		block, err := client.GetNowBlock()
		if err != nil {
			return nil, 0, err
		}
		tx, err := tron.DecodeTronTransaction(txHash)
		if err != nil {
			return nil, 0, err
		}
		transaction = *tx
		confirms = block.BlockHeader.RawData.Number - transaction.Height
	default:
		return nil, 0, errors.New("unsupported chain type")
	}
	return &transaction, confirms, nil
}

// @title	根据地址获取交易
// @param	Self		*Tron				模块实例
// @param	chainType	enum.ChainType		链类型
// @param	address		string				地址
// @param	token		*string				币种
// @param	endTime		time.Time			结束时间
// @return	_			[]model.Transaction	交易信息
// @return	_			error				异常信息
func (Self *Chain) GetTransactionsByAddress(chainType enum.ChainType, address string, token *string, endTime time.Time) ([]model.Transaction, error) {
	var transactions []model.Transaction
	switch chainType {
	case enum.ChainType_Tron:
		tron, _ := GetTron()
		txs, err := tron.GetTronTransactionsByAddress(address, token, endTime)
		if err != nil {
			return nil, err
		}
		transactions = txs
	default:
		return nil, errors.New("unsupported chain type")
	}
	return transactions, nil
}
