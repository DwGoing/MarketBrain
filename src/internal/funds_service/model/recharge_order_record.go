package model

import (
	"errors"
	"time"

	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RechargeOrderRecord struct {
	Record
	ExternalIdentity string    `gorm:"column:EXTERNAL_IDENTITY" mapstructure:"EXTERNAL_IDENTITY"`
	ExternalData     []byte    `gorm:"column:EXTERNAL_DATA" mapstructure:"EXTERNAL_DATA"`
	CallbackUrl      string    `gorm:"column:CALLBACK_URL" mapstructure:"CALLBACK_URL"`
	ChainType        string    `gorm:"column:CHAIN_TYPE" mapstructure:"CHAIN_TYPE"`
	Amount           float64   `gorm:"column:AMOUNT" mapstructure:"AMOUNT"`
	WalletIndex      int64     `gorm:"column:WALLET_INDEX" mapstructure:"WALLET_INDEX"`
	WalletAddress    string    `gorm:"column:WALLET_ADDRESS" mapstructure:"WALLET_ADDRESS"`
	Status           string    `gorm:"column:STATUS" mapstructure:"STATUS"`
	ExpireAt         time.Time `gorm:"column:EXPIRE_AT" mapstructure:"EXPIRE_AT"`
	TxHash           *string   `gorm:"column:TX_HASH" mapstructure:"TX_HASH"`
}

func CreateRechargeOrderRecord(client *gorm.DB, record *RechargeOrderRecord) (*RechargeOrderRecord, error) {
	record.Id = uuid.NewString()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	record.Status = enum.RechargeStatus_UNPAID.String()
	err := client.Transaction(func(tx *gorm.DB) error {
		var count int64
		result := tx.Table("RECHARGE_ORDER").Where("`EXTERNAL_IDENTITY` = ?", record.ExternalIdentity).Count(&count)
		if result.Error != nil {
			return result.Error
		}
		if count > 0 {
			return errors.New("order existed")
		}
		result = tx.Table("RECHARGE_ORDER").Create(&record)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func DeleteRechargeOrderRecords(client *gorm.DB, opt DeleteOption) error {
	result := client.Table("RECHARGE_ORDER").Where(opt.Conditions, opt.ConditionsParameters...).Delete(&RechargeOrderRecord{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateRechargeOrderRecords(client *gorm.DB, opt UpdateOption) error {
	if opt.Values == nil {
		return nil
	}
	opt.Values["UPDATED_AT"] = time.Now()
	result := client.Table("RECHARGE_ORDER").Where(opt.Conditions, opt.ConditionsParameters...).Updates(opt.Values)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetRechargeOrderRecords(client *gorm.DB, opt GetOption) ([]RechargeOrderRecord, int64, error) {
	client = client.Table("RECHARGE_ORDER").Order("`CREATED_AT` DESC")
	if opt.Conditions != "" {
		client = client.Where(opt.Conditions, opt.ConditionsParameters...)
	}
	var total int64
	result := client.Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if opt.PageSize > 0 && opt.PageIndex > 0 {
		client.Limit(int(opt.PageSize)).Offset(int((opt.PageIndex - 1) * opt.PageSize))
	}
	var records []RechargeOrderRecord
	result = client.Find(&records)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return records, total, nil
}
