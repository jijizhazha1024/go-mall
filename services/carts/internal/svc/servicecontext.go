package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jijizhazha1024/go-mall/services/carts/internal/config"
	"jijizhazha1024/go-mall/services/carts/internal/model"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB // 新增数据库连接
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.User,     // 配置中的 MySQL 用户名
		c.Mysql.Password, // 配置中的 MySQL 密码
		c.Mysql.Host,     // 配置中的 MySQL 主机
		c.Mysql.Port,     // 配置中的 MySQL 端口
		c.Mysql.Database, // 配置中的 MySQL 数据库名称
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err // 如果数据库连接失败，返回错误
	}

	if err := db.AutoMigrate(&model.Cart{}, &model.User{}, &model.Product{}); err != nil {
		logx.Errorf("AutoMigrate failed: %v", err)
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
	}, nil
}
