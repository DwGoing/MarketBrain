package model

import (
	"time"

	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RechargeRecord struct {
	Record
	ExternalIdentity string    `gorm:"column:EXTERNAL_IDENTITY"`
	ExternalData     []byte    `gorm:"column:EXTERNAL_DATA"`
	CallbackUrl      string    `gorm:"column:CALLBACK_URL"`
	ChainType        string    `gorm:"column:CHAIN_TYPE"`
	Amount           float64   `gorm:"column:AMOUNT"`
	WalletIndex      int64     `gorm:"column:WALLET_INDEX"`
	WalletAddress    string    `gorm:"column:WALLET_ADDRESS"`
	Status           string    `gorm:"column:STATUS"`
	ExpireAt         time.Time `gorm:"column:EXPIRE_AT"`
}

func CreateRechargeRecord(client *gorm.DB, record *RechargeRecord) (*RechargeRecord, error) {
	record.Id = uuid.NewString()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	record.Status = enum.RechargeStatus_UNPAID.String()
	result := client.Table("RECHARGE_RECORD").Create(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func DeleteRechargeRecords(client *gorm.DB, opt DeleteOption) error {
	result := client.Table("RECHARGE_RECORD").Where(opt.Conditions, opt.ConditionsParameters...).Delete(&RechargeRecord{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateRechargeRecords(client *gorm.DB, opt UpdateOption) error {
	if opt.Values == nil {
		return nil
	}
	var record RechargeRecord
	mapstructure.Decode(opt.Values, &record)
	record.UpdatedAt = time.Now()
	result := client.Table("RECHARGE_RECORD").Where(opt.Conditions, opt.ConditionsParameters...).Updates(record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetRechargeRecords(client *gorm.DB, opt GetOption) ([]RechargeRecord, int64, error) {
	client = client.Table("RECHARGE_RECORD").Order("`CREATED_AT` DESC")
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
	var records []RechargeRecord
	result = client.Find(&records)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return records, total, nil
}
