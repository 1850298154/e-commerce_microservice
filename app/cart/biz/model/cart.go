package model

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// type Cart struct {
// 	gorm.Model
// 	UserId    uint32 `gorm:"type:int(11);not null:index:idx_user_id"`
// 	ProductId uint32 `gorm:"type:int(11);not null"`
// 	Quantity  int32  `gorm:"type:int(11);not null"`
// }

// func (Cart) TableName() string {
// 	return "cart"
// }

// func (c Cart) AddItem(ctx context.Context, db *gorm.DB, item *Cart) error {
// 	// 检查商品是否存在
// 	var row Cart
// 	err := db.WithContext(ctx).Model(&Cart{}).
// 		Where(&Cart{UserId: item.UserId, ProductId: item.ProductId}).
// 		First(&row).Error
// 	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 		return err
// 	}

// 	// 商品存在，更新商品数量(当前数量+购物车数量)
// 	if row.ID > 0 {
// 		return db.WithContext(ctx).Model(&Cart{}).
// 			Where(&Cart{UserId: item.UserId, ProductId: item.ProductId}).
// 			UpdateColumn("quantity", gorm.Expr("quantity+?", item.Quantity)).Error
// 	}

// 	// 商品不存在，添加商品到购物车
// 	return db.WithContext(ctx).Create(item).Error
// }

// func (c Cart) EmptyCart(ctx context.Context, db *gorm.DB, userId uint32) error {
// 	// 清空购物车
// 	if userId == 0 {
// 		return errors.New("userId is empty")
// 	}
// 	return db.WithContext(ctx).Delete(&Cart{}, "user_id = ?", userId).Error
// }

// func (c Cart) GetCartByUserId(ctx context.Context, db *gorm.DB, userId uint32) ([]*Cart, error) {
// 	var cart []*Cart
// 	if userId == 0 {
// 		return nil, errors.New("userId is empty")
// 	}
// 	// 获取购物车列表
// 	err := db.WithContext(ctx).Model(&Cart{}).Where("user_id =?", userId).Find(&cart).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cart, nil
// }

type Cart struct {
	UserId    uint32
	ProductId string
	Quantity  int32
}

func (c Cart) AddItem(ctx context.Context, rdb *redis.Client, item *Cart) error {
	// 检查商品是否存在
	key := fmt.Sprintf("cart:%d", item.UserId)
	_, err := rdb.HGet(ctx, key, item.ProductId).Result()
	if err == redis.Nil {
		// 商品不存在，添加商品到购物车
		return rdb.HSet(ctx, key, item.ProductId, item.Quantity).Err()
	}
	// 商品存在，更新商品数量(当前数量+购物车数量)
	return rdb.HIncrBy(ctx, key, item.ProductId, int64(item.Quantity)).Err()
}

func (c Cart) EmptyCart(ctx context.Context, rdb *redis.Client, userId uint32) error {
	// 清空购物车
	if userId == 0 {
		return errors.New("userId is empty")
	}
	return rdb.Del(ctx, fmt.Sprintf("cart:%d", userId)).Err()
}

func (c Cart) GetCartByUserId(ctx context.Context, rdb *redis.Client, userId uint32) ([]*Cart, error) {
	cart := make([]*Cart, 0)
	if userId == 0 {
		return nil, errors.New("userId is empty")
	}
	// 获取购物车列表
	key := fmt.Sprintf("cart:%d", userId)
	items, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	for k, v := range items {
		// 将字符串转换为 int
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		// 将 int 转换为 int32
		i32 := int32(i)

		cart = append(cart, &Cart{
			UserId:    userId,
			ProductId: k,
			Quantity:  i32,
		})
	}
	return cart, nil
}
