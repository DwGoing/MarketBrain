package model

type Configs struct {
	Mnemonic     string                 `mapstructure:"MNEMONIC" json:"mnemonic"`
	ChainConfigs map[string]ChainConfig `mapstructure:"CHAIN_CONFIGS" json:"chainConfigs"`
}

type ChainConfig struct {
	USDT   string   `json:"usdt"`
	Nodes  []string `json:"nodes"`
	ApiKey string   `json:"apiKey"`
}
