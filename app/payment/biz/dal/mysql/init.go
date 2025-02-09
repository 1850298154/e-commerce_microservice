package mysql

import (
	"fmt"
	"os"

	"2501YTC/app/payment/biz/model"
	"2501YTC/app/payment/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	// dsn := "root:root@tcp(127.0.0.1:3306)/payment?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))
	fmt.Println("dsn:", dsn)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	_ = DB.AutoMigrate(&model.PaymentLog{})
	if err != nil {
		panic(err)
	}
}
