package config_module

import (
	"errors"
	"reflect"

	"funds-system/pkg/shared"
	"funds-system/pkg/storage_module"

	"github.com/ahmetb/go-linq"
	"github.com/mitchellh/mapstructure"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewConfigModule
type ConfigModule struct {
	StorageModule *storage_module.StorageModule `singleton:""`

	configs *shared.Configs
}

/*
@title	构造函数
@param 	module 	*ConfigModule 	模块实例
@return _ 		*ConfigModule	模块实例
@return _ 		error 			异常信息
*/
func NewConfigModule(module *ConfigModule) (*ConfigModule, error) {
	_, err := module.Load(false)
	if err != nil {
		return nil, err
	}
	m := map[string]float64{}
	m["0x466DD1e48570FAA2E7f69B75139813e4F8EF75c2"] = 5
	err = module.Set("CollectThresholds", &m)
	if err != nil {
		return nil, err
	}
	return module, nil
}

type ConfigRecord struct {
	Key   string `gorm:"column:KEY"`
	Value any    `gorm:"column:VALUE;serializer:json"`
}

/*
@title 	加载配置
@param 	Self	*ConfigModule	模块实例
@param 	cache	bool			是否使用缓存
@return _	 	*shared.Configs	是否使用缓存
@return _	 	error			异常信息
*/
func (Self *ConfigModule) Load(cache bool) (*shared.Configs, error) {
	if Self.configs == nil || !cache {
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
		Self.configs = configs
	}
	return Self.configs, nil
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
	_, err = Self.Load(true)
	if err != nil {
		return err
	}
	return nil
}
