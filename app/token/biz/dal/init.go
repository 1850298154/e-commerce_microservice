package dal

import (
	"2501YTC/app/token/biz/dal/redis"
)

func Init() {
	redis.Init()
	// mysql.Init()
}
