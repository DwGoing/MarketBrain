package chain_module

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/DwGoing/funds-system/pkg/hd_wallet"

	"github.com/alibaba/ioc-golang/extension/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewChainModule
type ChainModule struct {
	Nodes *config.ConfigSlice `config:",chain.nodes"`

	decimals map[common.Address]uint8
}

/*
@title	构造函数
@param 	module 	*ChainModule 	模块实例
@return _ 		*ChainModule	模块实例
@return _ 		error 			异常信息
*/
func NewChainModule(module *ChainModule) (*ChainModule, error) {
	module.decimals = make(map[common.Address]uint8)
	return module, nil
}

/*
@title	获取Eth客户端
@param 	Self 	*ChainModule 		模块实例
@return _ 		*ethclient.Client 	Eth客户端实例
@return _ 		error 				异常信息
*/
func (Self *ChainModule) getClient() (*ethclient.Client, error) {
	nodes := Self.Nodes.Value()
	node := nodes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(nodes))]
	client, err := ethclient.Dial(node.(string))
	if err != nil {
		return nil, err
	}
	return client, nil
}

/*
@title 获取Token的小数位数
@param 	Self 	*ChainModule 	模块实例
@param 	token 	string		 	Token地址
@return _ 		uint8 			小数位数
@return _ 		error 			异常信息
*/
func (Self *ChainModule) GetDecimals(token string) (uint8, error) {
	commonToken := common.HexToAddress(token)
	if v, ok := Self.decimals[commonToken]; ok {
		return v, nil
	} else {
		var decimals uint8 = 18
		if commonToken != common.HexToAddress("0x0") {
			client, err := Self.getClient()
			if err != nil {
				return 0, err
			}
			ierc20, err := NewIERC20(commonToken, client)
			if err != nil {
				return 0, err
			}
			decimals, err = ierc20.Decimals(nil)
			if err != nil {
				return 0, err
			}
		}
		Self.decimals[commonToken] = decimals
		return decimals, nil
	}
}

/*
@title 获取HD钱包
@param 	Self 		*ChainModule 		模块实例
@param 	mnemonic 	string 				助记词
@param 	password 	string	 			密码
@return _ 			*hd_wallet.HDWallet HD钱包
@return _ 			error 				异常信息
*/
func (Self *ChainModule) GetHDWallet(mnemonic string, password string) (*hd_wallet.HDWallet, error) {
	return hd_wallet.FromMnemonic(mnemonic, password)
}

/*
@title 按Token的小数位数转
@param 	Self 	*ChainModule 	模块实例
@param 	token 	string		 	Token地址
@param 	value 	*big.Float	 	转换前
@return _ 		*big.Int 		转换后
@return _ 		error 			异常信息
*/
func (Self *ChainModule) ConvertValue(token string, value *big.Float) (*big.Int, error) {
	decimals, err := Self.GetDecimals(token)
	if err != nil {
		return nil, err
	}
	convertedValue, _ := new(big.Float).Mul(value, big.NewFloat(float64(math.Pow10(int(decimals))))).Int(new(big.Int))
	return convertedValue, nil
}

/*
@title 按Token的小数位数转换
@param 	Self 	*ChainModule 	模块实例
@param 	token 	string		 	Token地址
@param 	value 	*big.Int 	 	转换前
@return _ 		*big.Float		转换后
@return _ 		error 			异常信息
*/
func (Self *ChainModule) UnconvertValue(token string, value *big.Int) (*big.Float, error) {
	decimals, err := Self.GetDecimals(token)
	if err != nil {
		return nil, err
	}
	unconvertedValue := new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(float64(math.Pow10(int(decimals)))))
	return unconvertedValue, nil
}

/*
@title	获取某个地址的余额
@param 	Self 	*ChainModule 	模块实例
@param 	address string 			查询地址
@param 	token 	string 			查询Token
@return _ 		*big.Int 		余额
@return _ 		error 			异常信息
*/
func (Self *ChainModule) GetBalance(address string, token string) (*big.Int, error) {
	// 检查参数
	if len(address) <= 0 {
		return nil, errors.New("address empty")
	}
	commonAddress := common.HexToAddress(address)
	commonToken := common.HexToAddress(token)
	client, err := Self.getClient()
	if err != nil {
		return nil, err
	}
	defer client.Client().Close()
	var balance *big.Int
	if commonToken == common.HexToAddress("0x0") {
		blockNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			return nil, err
		}
		balance, err = client.BalanceAt(context.Background(), commonAddress, big.NewInt(int64(blockNumber)))
		if err != nil {
			return nil, err
		}
	} else {
		ierc20, err := NewIERC20(commonToken, client)
		if err != nil {
			return nil, err
		}
		balance, err = ierc20.BalanceOf(nil, commonAddress)
		if err != nil {
			return nil, err
		}
	}
	return balance, nil
}

/*
@title	获取Gas价格
@param 	Self 	*ChainModule 	模块实例
@return _ 		*big.Int 		价格
@return _ 		error 			异常信息
*/
func (Self *ChainModule) GetGasPrice() (*big.Int, error) {
	client, err := Self.getClient()
	if err != nil {
		return nil, err
	}
	defer client.Client().Close()
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

/*
@title	转账
@param 	Self 		*ChainModule 		模块实例
@param 	privateKey 	*ecdsa.PrivateKey 	私钥
@param 	to 			string 				接收地址
@param 	token 		string 				转账Token
@param 	amount 		*big.Int 			数量
@return _ 			error 				异常信息
*/
func (Self *ChainModule) Transfer(privateKey *ecdsa.PrivateKey, to string, token string, amount *big.Int) error {
	// 检查参数
	if len(to) <= 0 {
		return errors.New("to empty")
	}
	commonTo := common.HexToAddress(to)
	commonToken := common.HexToAddress(token)
	client, err := Self.getClient()
	if err != nil {
		return err
	}
	defer client.Client().Close()
	var signedTx *types.Transaction
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	if commonToken == common.HexToAddress("0x0") {
		// 构造ETH转账交易
		tx := types.NewTransaction(nonce, commonTo, amount, 21000, gasPrice, nil)
		// 签名交易
		signedTx, err = types.SignTx(tx, types.LatestSignerForChainID(chainId), privateKey)
		if err != nil {
			return err
		}
	} else {
		// 构造ERC20转账交易
		ierc20, err := NewIERC20(commonToken, client)
		if err != nil {
			return err
		}
		// 签名交易
		transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
		if err != nil {
			return err
		}
		transactOpts.NoSend = true
		transactOpts.Nonce = big.NewInt(int64(nonce))
		transactOpts.GasLimit = uint64(300000)
		transactOpts.GasPrice = gasPrice
		signedTx, err = ierc20.Transfer(transactOpts, commonTo, amount)
		if err != nil {
			return err
		}
	}
	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	return nil
}
