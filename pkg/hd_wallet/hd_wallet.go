package hd_wallet

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/tyler-smith/go-bip39"
)

type HDWallet struct {
	seed []byte
}

type Version [2][4]byte

var (
	Version_xprv Version = [2][4]byte{{0x04, 0x88, 0xad, 0xe4}, {0x04, 0x88, 0xb2, 0x1e}}
	Version_yprv Version = [2][4]byte{{0x04, 0x9d, 0x78, 0x78}, {0x04, 0x88, 0x7c, 0xb2}}
	Version_Yprv Version = [2][4]byte{{0x02, 0x95, 0xb0, 0x05}, {0x04, 0x88, 0xb4, 0x3f}}
	Version_zprv Version = [2][4]byte{{0x04, 0xb2, 0x43, 0x0c}, {0x04, 0x88, 0x47, 0x46}}
	Version_Zprv Version = [2][4]byte{{0x02, 0xaa, 0x7a, 0x99}, {0x04, 0x88, 0x7e, 0xd3}}
	Version_tprv Version = [2][4]byte{{0x04, 0x35, 0x83, 0x94}, {0x04, 0x88, 0x87, 0xcf}}
	Version_uprv Version = [2][4]byte{{0x04, 0x4a, 0x4e, 0x28}, {0x04, 0x88, 0x52, 0x62}}
	Version_Uprv Version = [2][4]byte{{0x02, 0x42, 0x85, 0xb5}, {0x04, 0x88, 0x89, 0xef}}
	Version_vprv Version = [2][4]byte{{0x04, 0x5f, 0x18, 0xbc}, {0x04, 0x88, 0x1c, 0xf6}}
	Version_Vprv Version = [2][4]byte{{0x02, 0x57, 0x50, 0x48}, {0x04, 0x88, 0x54, 0x83}}
)

type Currency int8

const (
	Currency_BTC_Legacy       Currency = 1
	Currency_BTC_SegWit       Currency = 2
	Currency_BTC_NativeSegWit Currency = 3
	Currency_ETH              Currency = 4
)

/*
@title 	从种子导入钱包
@param 	seed   			[]byte 			种子
@return _ 				*HDWallet 		钱包实例
@return _ 				error 			异常信息
*/
func FromSeed(seed []byte) (*HDWallet, error) {
	if len(seed) < 16 || len(seed) > 64 {
		return nil, errors.New("seed invaild")
	}
	return &HDWallet{
		seed: seed,
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
@param 	Self   		*HDWallet 				HDWallet实例
@param 	version 	[4]byte					私钥ID
@param 	path 		string					派生路径
@return _			*hdkeychain.ExtendedKey 子私钥
@return _ 			error 					异常信息
*/
func (Self *HDWallet) DerivePrivateKey(version [4]byte, path string) (*hdkeychain.ExtendedKey, error) {
	masterKey, err := hdkeychain.NewMaster(
		Self.seed,
		&chaincfg.Params{HDPrivateKeyID: version},
	)
	if err != nil {
		return nil, err
	}
	parsedPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}
	accountKey := masterKey
	for _, n := range parsedPath {
		accountKey, err = accountKey.Child(n)
		if err != nil {
			return nil, err
		}
	}
	return accountKey, nil
}

/*
@title 	获取子钱包
@param 	Self   		*HDWallet 	HDWallet实例
@param 	currency 	Currency 	版本
@param 	coin 		int64 		币种
@param 	index 		Currency 	钱包索引
@return _			*Account	Account实例
@return _ 			error 		异常信息
*/
func (Self *HDWallet) GetAccount(currency Currency, index int64) (*Account, error) {
	var (
		version Version
		path    string
	)
	switch currency {
	case Currency_BTC_Legacy:
		version = Version_xprv
		path = "m/44'/0'/0'/0/"
	case Currency_BTC_SegWit:
		version = Version_xprv
		path = "m/49'/0'/0'/0/"
	case Currency_BTC_NativeSegWit:
		version = Version_xprv
		path = "m/84'/0'/0'/0/"
	case Currency_ETH:
		version = Version_xprv
		path = "m/44'/60'/0'/0/"
	default:
		return nil, errors.New("unsupportted currency")
	}
	path = fmt.Sprintf("%s%d", path, index)
	privateKey, err := Self.DerivePrivateKey(version[0], path)
	if err != nil {
		return nil, err
	}
	btcecPrivateKey, err := privateKey.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return &Account{
		Currency:   currency,
		PrivateKey: btcecPrivateKey,
	}, nil
}
