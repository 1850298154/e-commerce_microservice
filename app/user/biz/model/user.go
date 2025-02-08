package model

import (
	"context"

	"2501YTC/app/user/errno"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string `gorm:"unique_index"`
	PasswordHashed string `gorm:"type:varchar(255);not null"`
	Role           Role   `gorm:"default:1"` // 0表示Admin，1表示User
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
	err = u.db.WithContext(u.ctx).Model(&User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		err = errno.UserNotExistErr(err)
		klog.Error(err)
	}
	return
}

func (u *UserQuery) GetUserById(id uint32) (user *User, err error) {
	err = u.db.WithContext(u.ctx).Model(&User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		err = errno.UserNotExistErr(err)
		klog.Error(err)
	}
	return
}

func (u *UserQuery) UpdateUser(user *User) (err error) {
	err = u.db.WithContext(u.ctx).Model(&User{}).Updates(user).Error
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
