package module

import (
	"github.com/alibaba/ioc-golang/extension/config"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:constructFunc=NewStorage
type Storage struct {
	RedisConnectionString *config.ConfigString `config:",storage.redis"`
	MysqlConnectionString *config.ConfigString `config:",storage.mysql"`
}

// @title	构造函数
// @param 	service *Storage 	模块实例
// @return _ 		*Storage 	模块实例
// @return _ 		error 		异常信息
func NewStorage(module *Storage) (*Storage, error) {
	return module, nil
}

// @title 	获取redis连接
// @param 	Self	*Storage 		模块实例
// @return 	_ 		*redis.Client 	Redis客户端
// @return 	_ 		error 			异常信息
func (Self *Storage) GetRedisClient() (*redis.Client, error) {
	opt, err := redis.ParseURL(Self.RedisConnectionString.Value())
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return client, nil
}

// @title 	获取mysql连接
// @param 	Self	*Storage 	模块实例
// @return _ 		*gorm.DB	Mysql连接
// @return _ 		error 		异常信息
func (Self *Storage) GetMysqlClient() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(Self.MysqlConnectionString.Value()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
