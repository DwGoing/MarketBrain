package model

type Configs struct {
	ExpireTime                int64                  `mapstructure:"EXPIRE_TIME" json:"expireTime"`
	Mnemonic                  string                 `mapstructure:"MNEMONIC" json:"mnemonic"`
	ChainConfigs              map[string]ChainConfig `mapstructure:"CHAIN_CONFIGS" json:"chainConfigs"`
	WalletCollectionThreshold float64                `mapstructure:"WALLET_COLLECT_THRESHOLD" json:"walletCollectionThreshold"`
	MinGasThreshold           float64                `mapstructure:"MIN_GAS_THRESHOLD" json:"minGasThreshold"`
	TransferGasAmount         float64                `mapstructure:"TRANSFER_GAS_AMOUNT" json:"transferGasAmount"`
}

type ChainConfig struct {
	USDT             string   `json:"usdt"`
	RpcNodes         []string `json:"rpcNodes"`
	HttpNodes        []string `json:"httpNodes"`
	ApiKeys          []string `json:"apiKeys"`
	CollectionTarget string   `json:"collectionTarget"`
}
