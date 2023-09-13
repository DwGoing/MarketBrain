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
	Currency   Currency
	PrivateKey *btcec.PrivateKey
}

/*
@title 	获取钱包地址
@param 	Self   	*Account 	Account实例
@return _ 		string 		钱包地址
*/
func (Self *Account) GetAddress() string {
	address := ""
	switch Self.Currency {
	case Currency_BTC_Legacy:
		bytes := btcutil.Hash160(Self.PrivateKey.PubKey().SerializeCompressed())
		address = base58.CheckEncode(bytes, 0x00)
	case Currency_BTC_SegWit:
		bytes := btcutil.Hash160(Self.PrivateKey.PubKey().SerializeCompressed())
		bytes = append([]byte{0x00, 0x14}, bytes...)
		bytes = btcutil.Hash160(bytes)
		address = base58.CheckEncode(bytes, 0x05)
	case Currency_BTC_NativeSegWit:
		bytes := btcutil.Hash160(Self.PrivateKey.PubKey().SerializeCompressed())
		converted, err := bech32.ConvertBits(bytes, 8, 5, true)
		if err != nil {
			break
		}
		combined := make([]byte, len(converted)+1)
		combined[0] = 0x00
		copy(combined[1:], converted)
		address, err = bech32.Encode("bc", combined)
		if err != nil {
			break
		}
	case Currency_ETH:
		address = crypto.PubkeyToAddress(Self.PrivateKey.ToECDSA().PublicKey).Hex()
	case Currency_TRON:
		ethAddress := crypto.PubkeyToAddress(Self.PrivateKey.ToECDSA().PublicKey)
		tronAddress := make([]byte, 0)
		tronAddress = append(tronAddress, byte(0x41))
		tronAddress = append(tronAddress, ethAddress.Bytes()...)
		address = tronCommon.EncodeCheck(tronAddress)
	}
	return address
}
