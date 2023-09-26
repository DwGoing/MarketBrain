package hd_wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/btcutil/bech32"
	"github.com/ethereum/go-ethereum/crypto"
	tronCommon "github.com/fbsobreira/gotron-sdk/pkg/common"
)

type Account struct {
	Index      uint32
	PrivateKey *ecdsa.PrivateKey
}

/*
@title 	获取钱包地址
@param 	Self   	*Account 	Account实例
@return _ 		string 		钱包地址
*/
func (Self *Account) GetAddress() string {
	address := crypto.PubkeyToAddress(Self.PrivateKey.PublicKey)
	return address.Hex()
}
