package dal

import (
	"2501YTC/app/cart/biz/dal/mysql"
	"2501YTC/app/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
