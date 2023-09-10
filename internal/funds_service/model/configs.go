package model

import "github.com/DwGoing/MarketBrain/pkg/enum"

type Configs struct {
	Mnemonic string  `mapstructure:"MNEMONIC" json:"Mnemonic"`
	Chains   []Chain `mapstructure:"CHAINS" json:"Chains"`
}

type Chain struct {
	Type  enum.ChainType `json:"Type"`
	Nodes []string       `json:"Nodes"`
}
