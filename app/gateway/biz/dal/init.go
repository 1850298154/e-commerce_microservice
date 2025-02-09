package dal

import (
	"2501YTC/app/gateway/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
