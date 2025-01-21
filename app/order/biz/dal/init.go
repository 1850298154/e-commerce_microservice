package dal

import (
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/dal/mq"
)

func Init() {
	mysql.Init()
	mq.Init()
}
