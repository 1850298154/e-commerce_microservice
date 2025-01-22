package dal

import (
	"2501YTC/app/checkout/biz/dal/mysql"
	"2501YTC/app/checkout/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
