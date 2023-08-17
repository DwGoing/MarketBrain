package hd_wallet

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestFromMnemonic(t *testing.T) {
	var tests = []struct {
		mnemonic string
		password string
		want     string
	}{
		{"", "", "FromMnemonic Error: mnemonic empty"},
		{"math absorb sweet shrimp time smoke net pulp carbon gorilla expand", "", "FromMnemonic Error: mnemonic invaild"},
		{
			"math absorb sweet shrimp time smoke net pulp carbon gorilla expand payment", "",
			"3ff6cbb9f0e124a3e5cdd58f9ee4e4d87cd2aab60c91e45ae53ebc7f8eb35003a160d2cc5eec710b8fc4e5327482f793efa595281e613db239c6262b35c19362",
		},
	}
	for _, test := range tests {
		hdWallet, err := FromMnemonic(test.mnemonic, test.password)
		if err == nil {
			if common.Bytes2Hex(hdWallet.seed) != test.want {
				t.Error("Seed Error")
			}
		} else {
			if err.Error() != test.want {
				t.Error(err)
			}
		}
	}
}

func TestFromSeed(t *testing.T) {
	var tests = []struct {
		seed []byte
		want string
	}{
		{[]byte{}, "seed invaild"},
		{[]byte{0x0, 0x0}, "seed invaild"},
		{
			[]byte{63, 246, 203, 185, 240, 225, 36, 163, 229, 205, 213, 143, 158, 228, 228, 216, 124, 210, 170, 182, 12, 145, 228, 90, 229, 62, 188, 127, 142, 179, 80, 3, 161, 96, 210, 204, 94, 236, 113, 11, 143, 196, 229, 50, 116, 130, 247, 147, 239, 165, 149, 40, 30, 97, 61, 178, 57, 198, 38, 43, 53, 193, 147, 98},
			"3ff6cbb9f0e124a3e5cdd58f9ee4e4d87cd2aab60c91e45ae53ebc7f8eb35003a160d2cc5eec710b8fc4e5327482f793efa595281e613db239c6262b35c19362",
		},
	}
	for _, test := range tests {
		hdWallet, err := FromSeed(test.seed)
		if err == nil {
			if common.Bytes2Hex(hdWallet.seed) != test.want {
				t.Error("Seed Error")
			}
		} else {
			if err.Error() != test.want {
				t.Error(err)
			}
		}
	}
}

func TestDerivePrivateKey(t *testing.T) {
	mnemonic := "math absorb sweet shrimp time smoke net pulp carbon gorilla expand payment"
	var tests = []struct {
		path string
		want string
	}{
		{"", "ambiguous path: use 'm/' prefix for absolute paths, or no leading '/' for relative ones"},
		{"m/44'/60'/0'/0", "xprvA258xoK46SQ67fJuzV56VGTbBzrA89Qu1LxZJeGti1KtyVxH71TpxAXw6LAU7o7vyFmAvAU9bzjG1H2nCv9DVg3uqLTB9MukbL5hgLizoq7"},
		{"m/44'/60'/0'/0/0", "xprvA3apYdt417TRrAdLnuf1pjDXZoXDT8Hdo3YRUTfBTQtLcu5i6yQcxe4FhNP538Yh3iouZqQh6Ar4VsNqiEKhCGx9mpzZdMdtxJhrubQoLHz"},
	}
	for _, test := range tests {
		hdWallet, _ := FromMnemonic(mnemonic, "")
		privateKey, err := hdWallet.DerivePrivateKey(Version_xprv[0], test.path)
		if err == nil {
			if privateKey.String() != test.want {
				t.Error("Address Error")
			}
		} else {
			if err.Error() != test.want {
				t.Error(err)
			}
		}
	}
}

func TestGetAccount(t *testing.T) {
	mnemonic := "math absorb sweet shrimp time smoke net pulp carbon gorilla expand payment"
	var tests = []struct {
		currency Currency
		index    int64
		want     string
	}{
		{0, 0, "unsupportted currency"},
		{Currency_BTC_Legacy, 0, "1EG4SZFvXrYjBpz3QQp51MFnUivKCVXko2"},
		{Currency_BTC_SegWit, 0, "3JgMsmq79Ku8LWip3VCvF1pPQgbTrtx1gG"},
		{Currency_BTC_NativeSegWit, 0, "bc1q96cqwy3w2q7qelhecz3l2wu5ddc3lzfy0z6p0r"},
		{Currency_ETH, 0, "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E"},
	}
	for _, test := range tests {
		hdWallet, _ := FromMnemonic(mnemonic, "")
		account, err := hdWallet.GetAccount(test.currency, test.index)
		if err == nil {
			if account.GetAddress() != test.want {
				t.Error("Address Error")
			}
		} else {
			if err.Error() != test.want {
				t.Error(err)
			}
		}
	}
}
