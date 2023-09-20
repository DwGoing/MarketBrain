package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TransferRecord struct {
	Record
	ChainType   string  `gorm:"column:CHAIN_TYPE" mapstructure:"CHAIN_TYPE"`
	Token       *string `gorm:"column:TOKEN" mapstructure:"TOKEN"`
	FromIndex   int64   `gorm:"column:FROM_INDEX" mapstructure:"FROM_INDEX"`
	FromAddress string  `gorm:"column:FROM_ADDRESS" mapstructure:"FROM_ADDRESS"`
	To          string  `gorm:"column:TO" mapstructure:"TO"`
	Amount      float64 `gorm:"column:AMOUNT" mapstructure:"AMOUNT"`
	Status      string  `gorm:"column:STATUS" mapstructure:"STATUS"`
	Error       string  `gorm:"column:ERROR" mapstructure:"ERROR"`
	Remarks     string  `gorm:"column:REMARKS" mapstructure:"REMARKS"`
}

func CreateTransferRecord(client *gorm.DB, record *TransferRecord) (*TransferRecord, error) {
	record.Id = uuid.NewString()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	result := client.Table("`TRANSFER`").Create(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func DeleteTransferRecords(client *gorm.DB, opt DeleteOption) error {
	result := client.Table("`TRANSFER`").Where(opt.Conditions, opt.ConditionsParameters...).Delete(&TransferRecord{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateTransferRecords(client *gorm.DB, opt UpdateOption) error {
	if opt.Values == nil {
		return nil
	}
	var record TransferRecord
	mapstructure.Decode(opt.Values, &record)
	record.UpdatedAt = time.Now()
	result := client.Table("`TRANSFER`").Where(opt.Conditions, opt.ConditionsParameters...).Updates(record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetTransferRecords(client *gorm.DB, opt GetOption) ([]TransferRecord, int64, error) {
	client = client.Table("`TRANSFER`").Order("`CREATED_AT` DESC")
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
	var records []TransferRecord
	result = client.Find(&records)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return records, total, nil
}
