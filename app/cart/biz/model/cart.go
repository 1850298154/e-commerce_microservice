package model

import (
	"context"
	"errors"
	"fmt"
	"log"
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

// CartConfig 购物车配置
type CartConfig struct {
	MaxProductTypes int64
}

// CartService 购物车服务
type CartService struct {
	rdb    *redis.Client
	config CartConfig
}

// NewCartService 创建一个新的购物车服务实例
func NewCartService(rdb *redis.Client, config CartConfig) *CartService {
	return &CartService{
		rdb:    rdb,
		config: config,
	}
}

var cartServiceInstance *CartService
var cartConfigInstance CartConfig

// GetCartConfig 获取单例的 CartConfig 实例
func GetCartConfig() CartConfig {
	if cartConfigInstance.MaxProductTypes == 0 {
		cartConfigInstance = CartConfig{
			MaxProductTypes: 50,
		}
	}
	return cartConfigInstance
}

// GetCartService 获取单例的 CartService 实例
func GetCartService(rdb *redis.Client) *CartService {
	if cartServiceInstance == nil {
		config := GetCartConfig()
		cartServiceInstance = &CartService{
			rdb:    rdb,
			config: config,
		}
	}
	return cartServiceInstance
}

// getCartKey 生成购物车的 Redis 键
func (cs *CartService) getCartKey(userId uint32) string {
	return fmt.Sprintf("cart:%d", userId)
}

func (cs *CartService) AddItem(ctx context.Context, item *Cart) error {
	key := cs.getCartKey(item.UserId)
	// 检查购物车中商品的种类数量
	productCount, err := cs.rdb.HLen(ctx, key).Result()
	if err != nil {
		log.Printf("获取购物车商品种类数量失败，用户 ID: %d, 错误: %v", item.UserId, err)
		return err
	}
	if productCount >= cs.config.MaxProductTypes {
		// 检查商品是否已存在，如果存在则允许更新数量
		exists, err := cs.rdb.HExists(ctx, key, item.ProductId).Result()
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("购物车中商品种类已达到最大限制")
		}
	}
	_, err = cs.rdb.HGet(ctx, key, item.ProductId).Result()
	if err == redis.Nil {
		// 商品不存在，添加商品到购物车
		if err := cs.rdb.HSet(ctx, key, item.ProductId, item.Quantity).Err(); err != nil {
			log.Printf("添加商品到购物车失败，用户 ID: %d, 商品 ID: %s, 错误: %v", item.UserId, item.ProductId, err)
			return err
		}
		return nil
	}
	// 商品存在，更新商品数量(当前数量+购物车数量)
	if err := cs.rdb.HIncrBy(ctx, key, item.ProductId, int64(item.Quantity)).Err(); err != nil {
		log.Printf("更新购物车商品数量失败，用户 ID: %d, 商品 ID: %s, 错误: %v", item.UserId, item.ProductId, err)
		return err
	}
	return nil
}

func (cs *CartService) EmptyCart(ctx context.Context, userId uint32) error {
	// 清空购物车
	if userId == 0 {
		return errors.New("userId is empty")
	}
	key := cs.getCartKey(userId)
	if err := cs.rdb.Del(ctx, key).Err(); err != nil {
		log.Printf("清空购物车失败，用户 ID: %d, 错误: %v", userId, err)
		return err
	}
	return nil
}

func (cs *CartService) GetCartByUserId(ctx context.Context, userId uint32) ([]*Cart, error) {
	if userId == 0 {
		return nil, errors.New("userId is empty")
	}
	// 获取购物车列表
	key := cs.getCartKey(userId)
	items, err := cs.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("获取购物车信息失败，用户 ID: %d, 错误: %v", userId, err)
		return nil, err
	}
	cart := make([]*Cart, 0, len(items))
	for k, v := range items {
		// 将字符串转换为 int
		quantity, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, err
		}
		cart = append(cart, &Cart{
			UserId:    userId,
			ProductId: k,
			Quantity:  int32(quantity),
		})
	}
	return cart, nil
}
