package dal

import (
	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
