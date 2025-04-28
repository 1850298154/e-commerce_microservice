package model

import (
	"2501YTC/app/user/errno"
	"context"
	"encoding/json"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"time"
)

type UserLoginLog struct {
	gorm.Model
	UserID      uint32    `gorm:"not null;index" json:"user_id"`               // 用户ID，外键
	LoginTime   time.Time `gorm:"not null;index" json:"login_time"`            // 登录时间
	IPAddress   string    `gorm:"type:varchar(45);not null" json:"ip_address"` // 登录IP地址
	DeviceType  string    `gorm:"type:varchar(50)" json:"device_type"`         // 设备类型
	LoginStatus int8      `gorm:"not null" json:"login_status"`                // 登录状态（0失败，1成功）
	Location    string    `gorm:"type:varchar(255)" json:"location"`           // 地理位置
}

func (UserLoginLog) TableName() string { return "user_login_log" }

func (l *UserLoginLog) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}

func FromJSON(data []byte) (*UserLoginLog, error) {
	var log UserLoginLog
	err := json.Unmarshal(data, &log)
	return &log, err
}

type UserLoginLogQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserLoginLogQuery(ctx context.Context, db *gorm.DB) *UserLoginLogQuery {
	return &UserLoginLogQuery{ctx: ctx, db: db}
}

func (q *UserLoginLogQuery) Insert(log *UserLoginLog) (uint, error) {
	err := q.db.WithContext(q.ctx).Create(log).Error
	if err != nil {
		err = errno.CreateUserLoginLogErr(err)
		klog.Error(err)
		return 0, err
	}
	return 0, nil
}
