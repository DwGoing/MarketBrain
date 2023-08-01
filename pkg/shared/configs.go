package shared

/*
@title 配置表
*/
type Configs struct {
	Mnemonic          string             `mapstructure:"MNEMONIC"`
	WalletMaxNumber   int64              `mapstructure:"WALLET_MAX_NUMBER"`
	ExpireTime        int64              `mapstructure:"EXPIRE_TIME"`
	ExpireDelay       int64              `mapstructure:"EXPIRE_DELAY"`
	CollectThresholds map[string]float64 `mapstructure:"COLLECT_THRESHOLDS"`
}
