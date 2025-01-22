package model

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserId    uint32 `gorm:"type:int(11);not null:index:idx_user_id"`
	ProductId uint32 `gorm:"type:int(11);not null"`
	Quantity  uint32 `gorm:"type:int(11);not null"`
}

func (Cart) TableName() string {
	return "cart"
}

func AddItem(ctx context.Context, db *gorm.DB, item *Cart) error {
	// 检查商品是否存在
	var row Cart
	fmt.Printf("%+v", item)
	err := db.WithContext(ctx).Model(&Cart{}).
		Where(&Cart{UserId: item.UserId, ProductId: item.ProductId}).
		First(&row).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 商品存在，更新商品数量(当前数量+购物车数量)
	if row.ID > 0 {
		return db.WithContext(ctx).Model(&Cart{}).
			Where(&Cart{UserId: item.UserId, ProductId: item.ProductId}).
			UpdateColumn("quantity", gorm.Expr("quantity+?", item.Quantity)).Error
	}

	// 商品不存在，添加商品到购物车
	return db.WithContext(ctx).Create(item).Error
}
