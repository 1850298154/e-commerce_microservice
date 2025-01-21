package dal

import (
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
