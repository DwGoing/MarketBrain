package storage_module

import (
	"github.com/alibaba/ioc-golang/extension/config"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewStorageModule
// @title 		存储模块
// @description 包括缓存、持久化存储等
type StorageModule struct {
	RedisConnectionString *config.ConfigString `config:",storage.redis"`
	MysqlConnectionString *config.ConfigString `config:",storage.mysql"`
}

/*
@title	构造函数
@param 	module 	*StorageModule 	模块实例
@return _ 		*StorageModule 	模块实例
@return _ 		error 			异常信息
*/
func NewStorageModule(module *StorageModule) (*StorageModule, error) {
	return module, nil
}

/*
@title 	获取redis连接
@param 	service 	*StorageModule 	模块实例
@return _ 			*redis.Client 	Redis客户端
@return _ 			error 			异常信息
*/
func (Self *StorageModule) GetRedisConnection() (*redis.Client, error) {
	opt, err := redis.ParseURL(Self.RedisConnectionString.Value())
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return client, nil
}

/*
@title 	获取mysql连接
@param 	service 	*StorageModule 	模块实例
@return _ 			*gorm.DB	 	Mysql连接
@return _ 			error 			异常信息
*/
func (Self *StorageModule) GetMysqlConnection() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(Self.MysqlConnectionString.Value()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
