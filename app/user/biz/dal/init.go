package dal

import (
	"2501YTC/app/user/biz/dal/mysql"
	"2501YTC/app/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
