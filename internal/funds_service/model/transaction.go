package model

import (
	"github.com/DwGoing/MarketBrain/pkg/enum"
)

type Transaction struct {
	ChainType enum.ChainType
	Hash      string
	Height    int64
	TimeStamp int64
	Contract  *string
	From      string
	To        string
	Amount    float64
	Result    bool
}
