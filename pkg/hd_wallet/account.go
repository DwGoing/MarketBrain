package hd_wallet

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Index      uint32
	PrivateKey *ecdsa.PrivateKey
}

/*
@title 	获取钱包地址
@param 	Self   	*Wallet 	Wallet实例
@return _ 		string 		钱包地址
*/
func (Self *Account) GetAddress() string {
	address := crypto.PubkeyToAddress(Self.PrivateKey.PublicKey)
	return address.Hex()
}
