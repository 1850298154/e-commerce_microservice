package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       float32 `json:"price"`
	Categories  string  `json:"categories"`
}

type ProductsWithTotal struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

func (p Product) TableName() string {
	return "product"
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewProductQuery(ctx context.Context, db *gorm.DB) ProductQuery {
	return ProductQuery{ctx: ctx, db: db}
}

func (p *ProductQuery) Create(product Product) (Product, error) {
	err := p.db.WithContext(p.ctx).Create(&product).Error
	return product, err
}

func (p *ProductQuery) Update(id uint, product Product) (Product, error) {
	err := p.db.WithContext(p.ctx).Model(&Product{}).Where(&Product{Model: gorm.Model{ID: id}}).Updates(&product).Error
	return product, err
}

func (p *ProductQuery) Delete(productId int) error {
	return p.db.WithContext(p.ctx).Delete(&Product{}, productId).Error
}

func (p *ProductQuery) GetById(productId int) (product Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Where(&Product{Model: gorm.Model{ID: uint(productId)}}).First(&product).Error
	return
}

func (p *ProductQuery) GetAll() (products []Product, err error) {
	err = p.db.WithContext(p.ctx).Find(&products).Error
	return
}

func (p *ProductQuery) GetByCategory(category string, num int32, size int64) (products []Product, total int64, err error) {
	// 首先计算总数
	err = p.db.WithContext(p.ctx).Model(&Product{}).
		Where("categories like ?", "%"+category+"%").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 然后获取分页数据
	err = p.db.WithContext(p.ctx).
		Where("categories like ?", "%"+category+"%").
		Offset((int(num) - 1) * int(size)).
		Limit(int(size)).
		Find(&products).Error

	return products, total, err
}

func (p *ProductQuery) GetByName(name string, num int32, size int64, flag bool) (products []Product, total int64, err error) {
	// 首先计算总数
	var condition string
	if flag {
		// 模糊查询
		condition = "%" + name + "%"
	} else {
		// 详细查询
		condition = name
	}

	err = p.db.WithContext(p.ctx).Model(&Product{}).
		Where("name like ?", condition).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 然后获取分页数据
	err = p.db.WithContext(p.ctx).
		Where("name like ?", condition).
		Offset((int(num) - 1) * int(size)).
		Limit(int(size)).
		Find(&products).Error

	return products, total, err
}

type CachedProductQuery struct {
	productQuery ProductQuery
	cacheClient  *redis.Client
	prefix       string
}

func NewCachedProductQuery(pq ProductQuery, cacheClient *redis.Client) CachedProductQuery {
	return CachedProductQuery{productQuery: pq, cacheClient: cacheClient, prefix: "2501YTC"}
}

func (c CachedProductQuery) GetById(productId int) (product Product, err error) {
	cacheKey := fmt.Sprintf("%s:%s:%d", c.prefix, "product_by_id", productId)
	cachedResult := c.cacheClient.Get(c.productQuery.ctx, cacheKey)

	err = func() error {
		err1 := cachedResult.Err()
		if err1 != nil {
			return err1
		}
		cachedResultByte, err2 := cachedResult.Bytes()
		if err2 != nil {
			return err2
		}
		err3 := json.Unmarshal(cachedResultByte, &product)
		if err3 != nil {
			return err3
		}
		return nil
	}()
	if err != nil {
		product, err = c.productQuery.GetById(productId)
		if err != nil {
			return Product{}, err
		}
		encoded, err := json.Marshal(product)
		if err != nil {
			return product, nil
		}
		_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return
}

func (c CachedProductQuery) GetAll() (products []Product, err error) {
	cacheKey := fmt.Sprintf("%s:%s", c.prefix, "product_all")
	cachedResult := c.cacheClient.Get(c.productQuery.ctx, cacheKey)

	err = func() error {
		err1 := cachedResult.Err()
		if err1 != nil {
			return err1
		}
		cachedResultByte, err2 := cachedResult.Bytes()
		if err2 != nil {
			return err2
		}
		err3 := json.Unmarshal(cachedResultByte, &products)
		if err3 != nil {
			return err3
		}
		return nil
	}()
	if err != nil {
		products, err = c.productQuery.GetAll()
		if err != nil {
			return nil, err
		}
		encoded, err := json.Marshal(products)
		if err != nil {
			return products, nil
		}
		_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return
}

func (c CachedProductQuery) GetByCategory(category string, num int32, size int64) (products []Product, total int64, err error) {
	cacheKey := fmt.Sprintf("%s:%s:%s:%d-%d", c.prefix, "product_by_category", category, num, size)
	cachedResult := c.cacheClient.Get(c.productQuery.ctx, cacheKey)

	err = func() error {
		err1 := cachedResult.Err()
		if err1 != nil {
			return err1
		}
		cachedResultByte, err2 := cachedResult.Bytes()
		if err2 != nil {
			return err2
		}
		err3 := json.Unmarshal(cachedResultByte, &products)
		if err3 != nil {
			return err3
		}
		return nil
	}()
	if err != nil {
		products, total, err = c.productQuery.GetByCategory(category, num, size)
		if err != nil {
			return nil, 0, err
		}
		// 创建包含产品和总数的结构
		data := ProductsWithTotal{
			Products: products,
			Total:    total,
		}
		encoded, err := json.Marshal(data)
		if err != nil {
			return products, total, nil
		}
		_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return
}

func (c CachedProductQuery) Create(product Product) (Product, error) {
	product, err := c.productQuery.Create(product) //nolint
	if err != nil {
		return product, err
	}
	encoded, err := json.Marshal(product)
	if err != nil {
		return product, nil
	}
	cacheKey := fmt.Sprintf("%s:%s:%d", c.prefix, "product_by_id", product.ID)
	_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	_ = c.cacheClient.Del(c.productQuery.ctx, fmt.Sprintf("%s:%s", c.prefix, "product_all"))
	// 使用通配符构造匹配模式
	pattern := fmt.Sprintf("%s:%s:*", c.prefix, "product_by_category")
	// 查找所有匹配的键
	keys, err := c.cacheClient.Keys(c.productQuery.ctx, pattern).Result()
	if err != nil {
		return Product{}, fmt.Errorf("failed to fetch keys: %v", err)
	}

	if len(keys) == 0 {
		return product, nil
	}

	// 删除匹配的键
	_, err = c.cacheClient.Del(c.productQuery.ctx, keys...).Result()
	if err != nil {
		return Product{}, fmt.Errorf("failed to delete keys: %v", err)
	}
	return product, nil
}

func (c CachedProductQuery) Update(id uint, product Product) (Product, error) {
	product, err := c.productQuery.Update(id, product) //nolint
	if err != nil {
		return product, err
	}
	encoded, err := json.Marshal(product)
	if err != nil {
		return product, nil
	}
	cacheKey := fmt.Sprintf("%s:%s:%d", c.prefix, "product_by_id", product.ID)
	_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	_ = c.cacheClient.Del(c.productQuery.ctx, fmt.Sprintf("%s:%s", c.prefix, "product_all"))
	// 使用通配符构造匹配模式
	pattern := fmt.Sprintf("%s:%s:*", c.prefix, "product_by_category")
	// 查找所有匹配的键
	keys, err := c.cacheClient.Keys(c.productQuery.ctx, pattern).Result()
	if err != nil {
		return Product{}, fmt.Errorf("failed to fetch keys: %v", err)
	}

	if len(keys) == 0 {
		return product, nil
	}

	// 删除匹配的键
	_, err = c.cacheClient.Del(c.productQuery.ctx, keys...).Result()
	if err != nil {
		return Product{}, fmt.Errorf("failed to delete keys: %v", err)
	}
	return product, nil
}

func (c CachedProductQuery) Delete(productId int) error {
	err := c.productQuery.Delete(productId)
	if err != nil {
		return err
	}
	cacheKey := fmt.Sprintf("%s:%s:%d", c.prefix, "product_by_id", productId)
	_ = c.cacheClient.Del(c.productQuery.ctx, cacheKey)
	_ = c.cacheClient.Del(c.productQuery.ctx, fmt.Sprintf("%s:%s", c.prefix, "product_all"))
	// 使用通配符构造匹配模式
	pattern := fmt.Sprintf("%s:%s:*", c.prefix, "product_by_category")
	// 查找所有匹配的键
	keys, err := c.cacheClient.Keys(c.productQuery.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to fetch keys: %v", err)
	}

	if len(keys) == 0 {
		return nil
	}

	// 删除匹配的键
	_, err = c.cacheClient.Del(c.productQuery.ctx, keys...).Result()
	if err != nil {
		return fmt.Errorf("failed to delete keys: %v", err)
	}
	return nil
}
