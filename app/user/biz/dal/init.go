package dal

import (
	"2501YTC/app/user/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
