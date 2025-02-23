package model

import (
	"2501YTC/app/user/errno"
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string `gorm:"unique_index;type:varchar(255);not null"`
	PasswordHashed string `gorm:"type:varchar(255);not null"`
	Role           Role   `gorm:"default:1"`     // 1表示Admin，2表示User
	IsBanned       bool   `gorm:"default:false"` // 新增字段，默认不封禁
}

func (User) TableName() string {
	return "user"
}

type UserQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserQuery(ctx context.Context, db *gorm.DB) *UserQuery {
	return &UserQuery{ctx: ctx, db: db}
}

func (u *UserQuery) CreateUser(user *User) (id uint32, err error) {
	err = u.db.WithContext(u.ctx).Create(user).Error
	if err != nil {
		err = errno.CreateUserErr(err)
		klog.Error(err)
		return 0, err
	}
	return uint32(user.ID), nil
}

func (u *UserQuery) GetUserByEmail(email string) (user *User, err error) {
	err = u.db.WithContext(u.ctx).Model(&User{}).Where(&User{Email: email}).First(&user).Error
	if err != nil {
		err = errno.UserNotExistErr(err)
		klog.Error(err)
		return nil, err
	}
	return
}

func (u *UserQuery) GetUserById(id uint32) (user *User, err error) {
	err = u.db.WithContext(u.ctx).Model(&User{}).Where(&User{Model: gorm.Model{ID: uint(id)}}).First(&user).Error
	if err != nil {
		err = errno.UserNotExistErr(err)
		klog.Error(err)
		return nil, err
	}
	return
}

func (u *UserQuery) UpdateUser(user *User) (err error) {
	err = u.db.WithContext(u.ctx).Model(&User{}).Where(&User{Model: gorm.Model{ID: user.ID}}).Updates(user).Error
	if err != nil {
		err = errno.UpdateUserErr(err)
		klog.Error(err)
	}
	return
}

func (u *UserQuery) DeleteUser(id uint32) (err error) {
	err = u.db.WithContext(u.ctx).Model(&User{}).Delete(&User{}, id).Error
	if err != nil {
		err = errno.DeleteUserErr(err)
		klog.Error(err)
	}
	return
}
