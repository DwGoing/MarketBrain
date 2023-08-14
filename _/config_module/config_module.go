package config_module

import (
	"errors"
	"reflect"

	"github.com/DwGoing/OnlyPay/internal/module/storage_module"
	"github.com/DwGoing/OnlyPay/internal/shared"
	"github.com/ahmetb/go-linq"
	"github.com/mitchellh/mapstructure"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewConfigModule
type ConfigModule struct {
	StorageModule *storage_module.StorageModule `singleton:""`
}

/*
@title	构造函数
@param 	module 	*ConfigModule 	服务实例
@return _ 		*ConfigModule 	服务实例
@return _ 		error 			异常信息
*/
func NewConfigModule(module *ConfigModule) (*ConfigModule, error) {
	return module, nil
}

type ConfigRecord struct {
	Key   string `gorm:"column:KEY"`
	Value any    `gorm:"column:VALUE;serializer:json"`
}

/*
@title 	加载配置
@param 	Self	*ConfigModule	模块实例
@return _	 	*shared.Configs	配置表
@return _	 	error			异常信息
*/
func (Self *ConfigModule) Load() (*shared.Configs, error) {
	mysql, err := Self.StorageModule.GetMysqlConnection()
	if err != nil {
		return nil, err
	}
	sqlDB, err := mysql.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()
	var records []ConfigRecord
	result := mysql.Table("CONFIG").Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	configMap := map[string]any{}
	linq.From(records).ToMapByT(&configMap, func(i ConfigRecord) string {
		return i.Key
	}, func(i ConfigRecord) any {
		return i.Value
	})
	var configs *shared.Configs
	mapstructure.Decode(configMap, &configs)
	return configs, nil
}

/*
@title 	修改配置
@param 	Self	*ConfigModule	模块实例
@param 	key		string			Key
@param 	value	string			Value
@return _	 	error			异常信息
*/
func (Self *ConfigModule) Set(key string, value any) error {
	typeConfigs := reflect.TypeOf(shared.Configs{})
	field, ok := typeConfigs.FieldByName(key)
	if !ok {
		return errors.New("key invaild")
	}
	mysql, err := Self.StorageModule.GetMysqlConnection()
	if err != nil {
		return err
	}
	sqlDB, err := mysql.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	result := mysql.Table("CONFIG").Where("`KEY`=?", field.Tag.Get("mapstructure")).Select("VALUE").Updates(ConfigRecord{Value: value})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
