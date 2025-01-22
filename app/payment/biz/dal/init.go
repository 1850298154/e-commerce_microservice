package dal

import (
	"2501YTC/app/payment/biz/dal/mysql"
	"2501YTC/app/payment/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
