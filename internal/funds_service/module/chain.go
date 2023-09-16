package module

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	tronCommon "github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Chain struct {
	Config *Config `normal:""`
}

// @title	获取Tron客户端
// @param	Self	*Chain				模块实例
// @param	config	*ChainConfig		链配置
// @return	_		*client.GrpcClient	客户端
// @return	_		error				异常信息
func (Self *Chain) getTronClient(config ChainConfig) (*client.GrpcClient, error) {
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

// @title	解析交易
// @param	Self		*Chain			模块实例
// @param	chainType	enum.ChainType	链类型
// @param	txHash		string			交易Hash
// @return	_			bool			交易状态
// @return	_			string			合约地址
// @return	_			int64			时间戳
// @return	_			string			收款地址
// @return	_			float64			金额
// @return	_			int64			确认数
// @return	_			error			异常信息
func (Self *Chain) DecodeTransaction(chainType enum.ChainType, txHash string) (bool, string, int64, string, float64, int64, error) {
	var (
		result    bool
		address   string
		timeStamp int64
		to        string
		amount    float64
		confirms  int64
		err       error
	)
	config, err := Self.Config.Load()
	if err != nil {
		return result, address, timeStamp, to, amount, confirms, err
	}
	chainConfig, ok := config.ChainConfigs[enum.ChainType_TRON.String()]
	if !ok || len(chainConfig.Nodes) < 1 {
		err = errors.New("no chain config")
		goto finish
	}
	switch chainType {
	case enum.ChainType_TRON:
		client, err := Self.getTronClient(chainConfig)
		if err != nil {
			goto finish
		}
		tx, err := client.GetTransactionInfoByID(txHash)
		if err != nil {
			goto finish
		}
		result = tx.GetReceipt().GetResult() == core.Transaction_Result_SUCCESS
		address = tronCommon.EncodeCheck(tx.GetContractAddress())
		timeStamp = tx.GetBlockTimeStamp()
		lastestBlock, err := client.GetNowBlock()
		if err != nil {
			goto finish
		}
		confirms = lastestBlock.BlockHeader.RawData.Number - tx.BlockNumber
		log := tx.GetLog()[0]
		if tronCommon.BytesToHexString(log.GetTopics()[0]) != tronCommon.BytesToHexString(tronCommon.Keccak256([]byte("Transfer(address,address,uint256)"))) {
			return result, address, timeStamp, to, amount, confirms, errors.New("function not match")
		}
		to = tronCommon.EncodeCheck(append([]byte{0x41}, log.GetTopics()[2][12:]...))
		amount = float64(new(big.Int).SetBytes(log.Data).Uint64()) / 1e6
	default:
		err = errors.New("unsupported chain type")
	}
finish:
	return result, address, timeStamp, to, amount, confirms, err
}
