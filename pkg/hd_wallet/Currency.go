package hd_wallet

type Currency int8

const (
	Currency_BTC_Legacy       Currency = 1
	Currency_BTC_SegWit       Currency = 2
	Currency_BTC_NativeSegWit Currency = 3
	Currency_ETH              Currency = 4
)
