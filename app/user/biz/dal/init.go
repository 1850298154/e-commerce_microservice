package dal

import (
	"gomall/app/user/biz/dal/mysql"
	"gomall/app/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
