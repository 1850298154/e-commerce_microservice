package models

import (
	"context"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	UserID int32 `json:"user_id"`
	Role   int32 `json:"role"`
}

func (Token) TableName() string {
	return "token"
}

type TokenQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewTokenQuery(ctx context.Context, db *gorm.DB) *TokenQuery {
	return &TokenQuery{ctx: ctx, db: db}
}

func (q *TokenQuery) Create(token Token) (Token, error) {
	err := q.db.WithContext(q.ctx).Create(&token).Error
	return token, err
}

// GetByUserID 根据用户ID查找令牌
func (q *TokenQuery) GetByUserID(userID int32) (Token, error) {
	var token Token
	err := q.db.WithContext(q.ctx).Where("user_id = ?", userID).First(&token).Error
	return token, err
}

// Update 更新令牌信息
func (q *TokenQuery) Update(userID int32, token Token) (Token, error) {
	err := q.db.WithContext(q.ctx).Model(&Token{}).Where("user_id = ?", userID).Updates(&token).Error
	return token, err
}

// Delete 删除令牌记录
func (q *TokenQuery) Delete(userID int32) error {
	return q.db.WithContext(q.ctx).Where("user_id = ?", userID).Delete(&Token{}).Error
}
