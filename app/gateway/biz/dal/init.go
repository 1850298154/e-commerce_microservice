package dal

import (
	"2501YTC/app/gateway/biz/dal/mysql"
	"2501YTC/app/gateway/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
