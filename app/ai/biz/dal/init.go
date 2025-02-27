package dal

import (
	"2501YTC/app/ai/biz/dal/mysql"
	"2501YTC/app/ai/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
