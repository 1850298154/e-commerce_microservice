package mysql

import (
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"

	"2501YTC/app/user/biz/model"
	"2501YTC/app/user/conf"

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

	// 获取通用数据库对象 sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(conf.GetConf().MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.GetConf().MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.GetConf().MySQL.ConnMaxLifetime) * time.Second)

	if err = DB.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}
	klog.Infof("MySQL 初始化成功, DSN: %s", dsn)
	//if os.Getenv("GO_ENV") != "online" {
	//	needDemoData := !DB.Migrator().HasTable(&model.User{})
	//	DB.AutoMigrate( //nolint:errcheck
	//		&model.User{},
	//	)
	//	if needDemoData {
	//		DB.Exec("INSERT INTO `user` (`id`,`created_at`,`updated_at`,`email`,`password_hashed`,`role`) VALUES (1,'2023-12-26 09:46:19.852','2023-12-26 09:46:19.852','123@admin.com','$2a$10$jTvUFh7Z8Kw0hLV8WrAws.PRQTeuH4gopJ7ZMoiFvwhhz5Vw.bj7C', 0)")
	//	}
	//}
}
