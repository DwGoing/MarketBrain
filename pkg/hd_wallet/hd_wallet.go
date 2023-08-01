package hd_wallet

import (
	"crypto/ecdsa"
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/tyler-smith/go-bip39"
)

/*
@title HD钱包
*/
type HDWallet struct {
	masterKey *hdkeychain.ExtendedKey
}

/*
@title 	从种子导入钱包
@param 	seed   	[]byte 		种子
@return _ 		*HDWallet 	钱包实例
@return _ 		error 		异常信息
*/
func FromSeed(seed []byte) (*HDWallet, error) {
	if len(seed) <= 0 {
		return nil, errors.New("FromSeed Error: seed empty")
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return &HDWallet{
		masterKey: masterKey,
	}, nil
}

/*
@title 	从助记词导入钱包
@param 	mnemonic   	string 		助记词
@param 	password 	string 		密码
@return _ 			*HDWallet 	钱包实例
@return _ 			error 		异常信息
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
@param 	Self   	*HDWallet 				HDWallet实例
@param 	path 	accounts.DerivationPath 派生路径
@return _		*ecdsa.PrivateKey 		子私钥
@return _ 		error 					异常信息
*/
func (Self *HDWallet) derivePrivateKey(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	privateKey := Self.masterKey
	var err error
	for _, n := range path {
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
@return _		*Wallet		Wallet实例
@return _ 		error 		异常信息
*/
func (Self *HDWallet) GetWallet(index uint32) (*Wallet, error) {
	privateKey, err := Self.derivePrivateKey(append(accounts.DefaultRootDerivationPath, index))
	if err != nil {
		return nil, err
	}
	return &Wallet{
		Index:      index,
		PrivateKey: privateKey,
	}, nil
}
