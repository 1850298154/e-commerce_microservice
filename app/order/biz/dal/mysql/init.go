package mysql

import (
	"fmt"
	"time"

	"2501YTC/app/order/biz/model"
	"2501YTC/app/order/conf"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 初始化MySQL
func Init() {
	// 连接数据库
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, conf.GetConf().MySQL.User, conf.GetConf().MySQL.Password, conf.GetConf().MySQL.Host, conf.GetConf().MySQL.Port, conf.GetConf().MySQL.DBName)

	db, err := gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	DB = db

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(conf.GetConf().MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.GetConf().MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.GetConf().MySQL.ConnMaxLifetime) * time.Second)

	// 自动迁移
	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		panic(err)
	}
	klog.Infof("MySQL 初始化成功, DSN: %s", dsn)
}
