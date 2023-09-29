package model

import "github.com/DwGoing/MarketBrain/pkg/enum"

type Block struct {
	ChainType enum.ChainType
	Height    int64
	TimeStamp int64
}
