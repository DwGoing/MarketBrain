package module

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/DwGoing/MarketBrain/pkg/hd_wallet"
	"github.com/ahmetb/go-linq"
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
type Tron struct{}

// @title	获取Tron客户端
// @param	Self	*Tron				模块实例
// @return	_		*client.GrpcClient	客户端
// @return	_		error				异常信息
func (Self *Tron) GetTronRpcClient() (*client.GrpcClient, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return nil, err
	}
	chain, ok := config.Chains[enum.ChainType_Tron.String()]
	if !ok {
		return nil, errors.New("no chain config")
	}
	chainTron := chain.(model.Chain_Tron)
	if len(chainTron.RpcNodes) < 1 || len(chainTron.ApiKeys) < 1 {
		return nil, errors.New("no chain config")
	}
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chainTron.RpcNodes))))
	if err != nil {
		return nil, err
	}
	grpcClient := client.NewGrpcClient(chainTron.RpcNodes[index.Int64()])
	index, err = rand.Int(rand.Reader, big.NewInt(int64(len(chainTron.ApiKeys))))
	if err != nil {
		return nil, err
	}
	err = grpcClient.SetAPIKey(chainTron.ApiKeys[index.Int64()])
	if err != nil {
		return nil, err
	}
	err = grpcClient.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return grpcClient, nil
}

// @title	获取Tron客户端
// @param	Self	*Tron			模块实例
// @return	_		*http.Client	客户端
// @return	_		error			异常信息
func (Self *Tron) GetTronHttpClient() (*http.Client, error) {
	configModule, _ := GetConfig()
	config, err := configModule.Load()
	if err != nil {
		return nil, err
	}
	chain, ok := config.Chains[enum.ChainType_Tron.String()]
	if !ok {
		return nil, errors.New("no chain config")
	}
	chainTron := chain.(model.Chain_Tron)
	if len(chainTron.HttpNodes) < 1 {
		return nil, errors.New("no chain config")
	}
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chainTron.RpcNodes))))
	if err != nil {
		return nil, err
	}
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.Dial(network, chainTron.HttpNodes[index.Int64()])
		},
	}
	return &http.Client{
		Transport: &transport,
	}, nil
}

// @title	获取当前高度
// @param	Self		*Tron			模块实例
// @return	_			int64			当前高度
// @return	_			error			异常信息
func (Self *Tron) GetCurrentHeight() (int64, error) {
	client, err := Self.GetTronRpcClient()
	if err != nil {
		return 0, err
	}
	block, err := client.GetNowBlock()
	if err != nil {
		return 0, err
	}
	return block.BlockHeader.RawData.Number, nil
}

// @title	获取钱包余额
// @param	Self		*Tron			模块实例
// @param	contract	*string			合约地址
// @param	address		string			钱包地址
// @return	_			float64			余额
// @return	_			error			异常信息
func (Self *Tron) GetBalance(contract *string, address string) (float64, error) {
	client, err := Self.GetTronRpcClient()
	if err != nil {
		return 0, err
	}
	var balance float64
	if contract == nil {
		account, err := client.GetAccount(address)
		if err != nil {
			return 0, err
		}
		balance, _ = new(big.Float).Quo(new(big.Float).SetInt64(account.Balance), big.NewFloat(1e6)).Float64()
	} else {
		balanceBigInt, err := client.TRC20ContractBalance(address, *contract)
		if err != nil {
			return 0, err
		}
		decimalsBigInt, err := client.TRC20GetDecimals(*contract)
		if err != nil {
			return 0, err
		}
		balance, _ = new(big.Float).Quo(new(big.Float).SetInt(balanceBigInt), big.NewFloat(math.Pow10(int(decimalsBigInt.Int64())))).Float64()
	}
	return balance, nil
}

// @title	发送Tron交易
// @param	Self		*Tron					模块实例
// @param	client		*client.GrpcClient		客户端
// @param	privateKey	*ecdsa.PrivateKey		私钥
// @param	tx			*core.Transaction		交易
// @param	waitReceipt	*client.GrpcClient		是否等待结果
// @return	_			*core.TransactionInfo	交易信息
// @return	_			error					异常信息
func (Self *Tron) SendTronTransaction(client *client.GrpcClient, privateKey *ecdsa.PrivateKey, tx *core.Transaction, waitReceipt bool) (*core.TransactionInfo, error) {
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

// @title	发送代币
// @param	Self		*Tron				模块实例
// @param	token		*string				合约地址
// @param	from		*hd_wallet.Account	发送账户
// @param	to			string				接收地址
// @param	amount		float64				数量
// @return	_			error				异常信息
func (Self *Tron) Transfer(token *string, from *hd_wallet.Account, to string, amount float64) (string, error) {
	client, err := Self.GetTronRpcClient()
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
	txInfo, err := Self.SendTronTransaction(client, from.PrivateKey.ToECDSA(), tx.Transaction, true)
	if err != nil {
		return "", err
	}
	return common.Bytes2Hex(txInfo.GetId()), nil
}

// @title	查询块信息
// @param	Self		*Chain				模块实例
// @param	number		int64				块高
// @return	_			*model.Block		块信息
// @return	_			error				异常信息
func (Self *Tron) GetBlockByNumber(number int64) (*model.Block, error) {
	result := model.Block{}
	client, err := Self.GetTronRpcClient()
	if err != nil {
		return nil, err
	}
	block, err := client.GetBlockByNum(number)
	if err != nil {
		return nil, err
	}
	result.ChainType = enum.ChainType_Tron
	result.Height = block.BlockHeader.RawData.Number
	result.TimeStamp = block.BlockHeader.RawData.Timestamp
	return &result, nil
}

// @title	解析交易
// @param	Self		*Tron				模块实例
// @param	txHash		string				交易Hash
// @return	_			*model.Transaction	交易信息
// @return	_			error				异常信息
func (Self *Tron) DecodeTronTransaction(txHash string) (*model.Transaction, error) {
	client, err := Self.GetTronRpcClient()
	if err != nil {
		return nil, err
	}
	result := model.Transaction{
		ChainType: enum.ChainType_Tron,
		Hash:      txHash,
	}
	tx, err := client.GetTransactionInfoByID(txHash)
	if err != nil {
		return nil, err
	}
	result.Height = tx.BlockNumber
	result.TimeStamp = tx.GetBlockTimeStamp()
	coreTx, err := client.GetTransactionByID(txHash)
	if err != nil {
		return nil, err
	}
	contracts := coreTx.RawData.GetContract()
	if len(contracts) < 1 {
		return nil, errors.New("not transfer transaction")
	}
	if contracts[0].Type != core.Transaction_Contract_TransferContract &&
		contracts[0].Type != core.Transaction_Contract_TriggerSmartContract {
		return nil, errors.New("not transfer transaction")
	}
	if tx.ContractAddress == nil {
		var contract core.TransferContract
		err = contracts[0].GetParameter().UnmarshalTo(&contract)
		if err != nil {
			return nil, err
		}
		result.From = common.EncodeCheck(contract.GetOwnerAddress())
		result.To = common.EncodeCheck(contract.GetToAddress())
		result.Amount, _ = new(big.Float).Quo(new(big.Float).SetInt64(contract.GetAmount()), big.NewFloat(1e6)).Float64()
	} else {
		contractAddress := common.EncodeCheck(tx.ContractAddress)
		result.Contract = &contractAddress
		var contract core.TriggerSmartContract
		err = contracts[0].GetParameter().UnmarshalTo(&contract)
		if err != nil {
			return nil, err
		}
		logs := tx.GetLog()
		if len(logs) < 1 {
			return nil, errors.New("not transfer transaction")
		}
		log := logs[0]
		topics := log.GetTopics()
		if len(topics) < 3 {
			return nil, errors.New("not transfer transaction")
		}
		// 签名校验
		if common.BytesToHexString(topics[0]) != common.BytesToHexString(common.Keccak256([]byte("Transfer(address,address,uint256)"))) {
			return nil, errors.New("not transfer transaction")
		}
		result.From = common.EncodeCheck(contract.GetOwnerAddress())
		result.To = common.EncodeCheck(append([]byte{0x41}, topics[2][12:]...))
		decimalsBigInt, err := client.TRC20GetDecimals(contractAddress)
		if err != nil {
			return nil, err
		}
		result.Amount, _ = new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).SetBytes(log.Data)), big.NewFloat(math.Pow10(int(decimalsBigInt.Int64())))).Float64()
	}
	receiptResult := tx.GetReceipt().GetResult()
	result.Result = receiptResult == core.Transaction_Result_SUCCESS
	return &result, nil
}

// @title	从块中获取交易
// @param	Self		*Tron					模块实例
// @param	start		int64					开始高度
// @param	end			int64					结束高度
// @return	_			[]model.Transaction		交易信息
// @return	_			error					异常信息
func (Self *Tron) GetTronTransactionsFromBlocks(start int64, end int64) ([]model.Transaction, error) {
	client, err := Self.GetTronRpcClient()
	if err != nil {
		return nil, err
	}
	blocklist, err := client.GetBlockByLimitNext(start, end)
	if err != nil {
		return nil, err
	}
	blocks := blocklist.GetBlock()
	result := []model.Transaction{}
	for _, block := range blocks {
		transactions := block.GetTransactions()
		for _, transaction := range transactions {
			tx, err := Self.DecodeTronTransaction(common.Bytes2Hex(transaction.GetTxid()))
			if err != nil {
				continue
			}
			result = append(result, *tx)
		}
	}
	return result, nil
}

type GetTronTransactionsByAddressResponse struct {
	Data []GetTronTransactionsByAddressResponse_Trc20Transaction `json:"data"`
}

type GetTronTransactionsByAddressResponse_Trc20Transaction struct {
	TransactionId  string                                                          `json:"transaction_id"`
	BlockTimestamp int64                                                           `json:"block_timestamp"`
	From           string                                                          `json:"from"`
	To             string                                                          `json:"to"`
	Value          string                                                          `json:"value"`
	TokenInfo      GetTronTransactionsByAddressResponse_Trc20Transaction_TokenInfo `json:"token_info"`
}

type GetTronTransactionsByAddressResponse_Trc20Transaction_TokenInfo struct {
	Address  string `json:"address"`
	Decimals int64  `json:"decimals"`
}

// @title	根据地址获取交易
// @param	Self		*Tron				模块实例
// @param	address		string				地址
// @param	token		*string				币种
// @param	endTime		time.Time			结束时间
// @return	_			[]model.Transaction	交易信息
// @return	_			error				异常信息
func (Self *Tron) GetTronTransactionsByAddress(address string, token *string, endTime time.Time) ([]model.Transaction, error) {
	client, err := Self.GetTronHttpClient()
	if err != nil {
		return nil, err
	}
	var transactions []model.Transaction
	if token == nil {
		// 未实现
	} else {
		url := fmt.Sprintf("/v1/accounts/%s/transactions/trc20?only_confirmed=true&contract_address=%s&min_timestamp=%d",
			address, *token, endTime.UnixMilli(),
		)
		response, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var res GetTronTransactionsByAddressResponse
		err = json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}
		linq.From(res.Data).SelectT(func(item GetTronTransactionsByAddressResponse_Trc20Transaction) model.Transaction {
			amountBitInt, _ := new(big.Int).SetString(item.Value, 10)
			amount, _ := new(big.Float).Quo(new(big.Float).SetInt(amountBitInt), big.NewFloat(math.Pow10(int(item.TokenInfo.Decimals)))).Float64()
			return model.Transaction{
				ChainType: enum.ChainType_Tron,
				Hash:      item.TransactionId,
				TimeStamp: item.BlockTimestamp,
				Contract:  &item.TokenInfo.Address,
				From:      item.From,
				To:        item.To,
				Amount:    amount,
				Result:    true,
			}
		}).ToSlice(&transactions)
	}
	return transactions, nil
}
