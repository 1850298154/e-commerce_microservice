package dal

import (
	"2501YTC/app/order/biz/dal/mq"
	"2501YTC/app/order/biz/dal/mysql"
)

func Init() {
	mysql.Init()
	mq.Init()
}
