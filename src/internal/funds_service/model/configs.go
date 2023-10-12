package model

type Configs struct {
	Chains                    map[string]any     `mapstructure:"CHAINS" json:"chains"`
	ExpireTime                int64              `mapstructure:"EXPIRE_TIME" json:"expireTime"`
	Mnemonic                  string             `mapstructure:"MNEMONIC" json:"mnemonic"`
	WalletCollectionThreshold map[string]float64 `mapstructure:"WALLET_COLLECT_THRESHOLD" json:"walletCollectionThreshold"`
	MinGasThreshold           map[string]float64 `mapstructure:"MIN_GAS_THRESHOLD" json:"minGasThreshold"`
	TransferGasAmount         map[string]float64 `mapstructure:"TRANSFER_GAS_AMOUNT" json:"transferGasAmount"`
	CollectWallets            map[string]string  `mapstructure:"COLLECT_WALLETS" json:"collectWallets"`
	PaymentCurrencies         map[string]string  `mapstructure:"PAYMENT_CURRENCIES" json:"paymentCurrencies"`
}

type Chain_Tron struct {
	RpcNodes  []string `mapstructure:"rpcNodes" json:"rpcNodes"`
	HttpNodes []string `mapstructure:"httpNodes" json:"httpNodes"`
	ApiKeys   []string `mapstructure:"apiKeys" json:"apiKeys"`
}
