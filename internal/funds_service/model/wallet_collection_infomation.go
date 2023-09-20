package model

import (
	"github.com/DwGoing/MarketBrain/pkg/enum"
)

type WalletCollectionInfomation struct {
	Index     int64
	ChainType enum.ChainType
	Address   string
}
