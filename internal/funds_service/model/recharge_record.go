package model

import (
	"time"
)

// @title 充值状态
type RechargeStatus int8

const (
	RechargeStatus_UNPAID RechargeStatus = 1
	RechargeStatus_PAID   RechargeStatus = 2
)

func (e RechargeStatus) String() string {
	switch e {
	case RechargeStatus_UNPAID:
		return "UNPAID"
	case RechargeStatus_PAID:
		return "PAID"
	default:
		return "UNKNOWN"
	}
}

// @title 充值记录
type RechargeRecord struct {
	Record
	ExternalIdentity string         `gorm:"column:EXTERNAL_IDENTITY"`
	ExternalData     []byte         `gorm:"column:EXTERNAL_DATA"`
	CallbackUrl      string         `gorm:"column:CALLBACK_URL"`
	Token            string         `gorm:"column:TOKEN"`
	Amount           string         `gorm:"column:AMOUNT"`
	WalletIndex      int64          `gorm:"column:WALLET_INDEX"`
	WalletAddress    string         `gorm:"column:WALLET_ADDRESS"`
	BeforeBalance    string         `gorm:"column:BDFORE_BALANCE"`
	AfterBalance     string         `gorm:"column:AFTER_BALANCE"`
	Status           RechargeStatus `gorm:"column:STATUS"`
	ExpireAt         time.Time      `gorm:"column:EXPIRE_AT"`
}
