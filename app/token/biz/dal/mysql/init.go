package mysql

import (
	"fmt"

	models "2501YTC/app/token/biz/model"
	"2501YTC/app/token/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, conf.GetConf().MySQL.User, conf.GetConf().MySQL.Password, conf.GetConf().MySQL.Host, conf.GetConf().MySQL.Port, conf.GetConf().MySQL.DBName)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	if err = DB.AutoMigrate(&models.Token{}); err != nil {
		panic(fmt.Sprintf("AutoMigrate failed: %v", err))
	}
}
