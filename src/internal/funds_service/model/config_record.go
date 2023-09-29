package model

import (
	"gorm.io/gorm"
)

type ConfigRecord struct {
	Key   string `gorm:"column:KEY"`
	Value any    `gorm:"column:VALUE;serializer:json"`
}

func UpdateConfigRecords(client *gorm.DB, configs []ConfigRecord) error {
	err := client.Transaction(func(tx *gorm.DB) error {
		for _, config := range configs {
			result := tx.Table("CONFIG").Where("`KEY`=?", config.Key).
				Updates(ConfigRecord{Value: config.Value})
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func GetConfigRecords(client *gorm.DB) ([]ConfigRecord, error) {
	var records []ConfigRecord
	result := client.Table("CONFIG").Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}
