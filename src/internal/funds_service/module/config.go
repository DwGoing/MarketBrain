package module

import (
	"github.com/DwGoing/MarketBrain/internal/funds_service/model"
	"github.com/ahmetb/go-linq"
	"github.com/mitchellh/mapstructure"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Config struct{}

// @title	更新配置
// @param	Self	*Config		模块实例
// @return	_		*Configs	配置
// @return	_		error		异常信息
func (Self *Config) Set(configs map[string]any) error {
	configRecords := []model.ConfigRecord{}
	for k, v := range configs {
		configRecords = append(configRecords, model.ConfigRecord{
			Key:   k,
			Value: v,
		})
	}
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = model.UpdateConfigRecords(mysqlClient, configRecords)
	if err != nil {
		return err
	}
	return nil
}

// @title	加载配置
// @param	Self	*Config		模块实例
// @return	_		*Configs	配置
// @return	_		error		异常信息
func (Self *Config) Load() (*model.Configs, error) {
	storageModule, _ := GetStorage()
	mysqlClient, err := storageModule.GetMysqlClient()
	if err != nil {
		return nil, err
	}
	db, err := mysqlClient.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	configRecords, err := model.GetConfigRecords(mysqlClient)
	if err != nil {
		return nil, err
	}
	configMap := map[string]any{}
	linq.From(configRecords).ToMapByT(&configMap, func(item model.ConfigRecord) string {
		return item.Key
	}, func(item model.ConfigRecord) any {
		return item.Value
	})
	var configs *model.Configs
	mapstructure.Decode(configMap, &configs)
	return configs, nil
}
