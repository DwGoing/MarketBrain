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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
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
// @param	nodes 	[]string			链配置
// @param	apiKey 	string				ApiKey
// @return	_		*client.GrpcClient	客户端
// @return	_		error				异常信息
func (Self *Tron) GetTronClient(nodes []string, apiKey string) (*client.GrpcClient, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(nodes))))
	if err != nil {
		return nil, err
	}
	grpcClient := client.NewGrpcClient(nodes[index.Int64()])
	err = grpcClient.SetAPIKey(apiKey)
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

// @title	解析交易
// @param	Self		*Tron				模块实例
// @param	txHash		string				交易Hash
// @return	_			*model.Transaction	交易信息
// @return	_			error				异常信息
func (Self *Tron) DecodeTransaction(client *client.GrpcClient, txHash string) (*model.Transaction, error) {
	result := model.Transaction{
		ChainType: enum.ChainType_TRON,
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
// @param	client		*client.GrpcClient		客户端
// @param	start		int64					开始高度
// @param	end			int64					结束高度
// @return	_			[]model.Transaction		交易信息
// @return	_			error					异常信息
func (Self *Tron) GetTransactionsFromBlocks(client *client.GrpcClient, start int64, end int64) ([]model.Transaction, error) {
	blocklist, err := client.GetBlockByLimitNext(start, end)
	if err != nil {
		return nil, err
	}
	blocks := blocklist.GetBlock()
	result := []model.Transaction{}
	for _, block := range blocks {
		transactions := block.GetTransactions()
		for _, transaction := range transactions {
			tx, err := Self.DecodeTransaction(client, common.Bytes2Hex(transaction.GetTxid()))
			if err != nil {
				continue
			}
			result = append(result, *tx)
		}
	}
	return result, nil
}
