package hd_wallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/tyler-smith/go-bip39"
)

type HDWallet struct {
	masterKey *hdkeychain.ExtendedKey
}

type Currency int64

const (
	Currency_ETH Currency = 60
)

/*
@title 	从种子导入钱包
@param 	seed   			[]byte 			种子
@return _ 				*HDWallet 		钱包实例
@return _ 				error 			异常信息
*/
func FromSeed(seed []byte) (*HDWallet, error) {
	if len(seed) <= 0 {
		return nil, errors.New("FromSeed Error: seed empty")
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.Params{HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}})
	if err != nil {
		return nil, err
	}
	return &HDWallet{
		masterKey: masterKey,
	}, nil
}

/*
@title 	从助记词导入钱包
@param 	mnemonic   		string 			助记词
@param 	password 		string 			密码
@return _ 				*HDWallet 		钱包实例
@return _ 				error 			异常信息
*/
func FromMnemonic(mnemonic string, password string) (*HDWallet, error) {
	if mnemonic == "" {
		return nil, errors.New("FromMnemonic Error: mnemonic empty")
	}
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("FromMnemonic Error: mnemonic invaild")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, err
	}
	return FromSeed(seed)
}

/*
@title 	派生子私钥
@param 	Self   	*HDWallet 			HDWallet实例
@param 	path 	string				派生路径
@return _		*ecdsa.PrivateKey 	子私钥
@return _ 		error 				异常信息
*/
func (Self *HDWallet) DerivePrivateKey(path string) (*ecdsa.PrivateKey, error) {
	parsedPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}
	privateKey := Self.masterKey
	for _, n := range parsedPath {
		privateKey, err = privateKey.Child(n)
		if err != nil {
			return nil, err
		}
	}
	btcecPrivateKey, err := privateKey.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return btcecPrivateKey.ToECDSA(), nil
}

/*
@title 	获取子钱包
@param 	Self   	*HDWallet 	HDWallet实例
@param 	index 	uint32 		钱包索引
@return _		*Account	Account实例
@return _ 		error 		异常信息
*/
func (Self *HDWallet) GetAccount(currency Currency, index uint32) (*Account, error) {
	privateKey, err := Self.DerivePrivateKey(fmt.Sprintf("m/44'/%d'/0'/0/%d", currency, index))
	if err != nil {
		return nil, err
	}
	return &Account{
		Index:      index,
		PrivateKey: privateKey,
	}, nil
}
