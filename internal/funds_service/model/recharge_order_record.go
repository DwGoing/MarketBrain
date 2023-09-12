package model

import (
	"time"

	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RechargeOrderRecord struct {
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

func CreateRechargeOrderRecord(client *gorm.DB, record *RechargeOrderRecord) (*RechargeOrderRecord, error) {
	record.Id = uuid.NewString()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	record.Status = enum.RechargeStatus_UNPAID.String()
	result := client.Table("RECHARGE_ORDER").Create(&record)
	if result.Error != nil {
		return nil, result.Error
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
	var record RechargeOrderRecord
	mapstructure.Decode(opt.Values, &record)
	record.UpdatedAt = time.Now()
	result := client.Table("RECHARGE_ORDER").Where(opt.Conditions, opt.ConditionsParameters...).Updates(record)
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
