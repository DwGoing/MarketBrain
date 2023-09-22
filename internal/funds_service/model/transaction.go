package model

import (
	"github.com/DwGoing/MarketBrain/pkg/enum"
)

type Transaction struct {
	ChainType enum.ChainType
	Contract  *string
	Hash      string
	TimeStamp int64
	From      string
	To        string
	Amount    float64
	Result    bool
}
