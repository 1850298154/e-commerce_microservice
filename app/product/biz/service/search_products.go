package service

import (
	"context"
	"fmt"
	"math"

	"github.com/cloudwego/kitex/pkg/klog"

	"2501YTC/app/product/biz/dal/meili"
	"2501YTC/app/product/utils/apiErr"

	"github.com/meilisearch/meilisearch-go"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

type SearchProductsService struct {
	ctx context.Context
} // NewSearchProductsService new SearchProductsService
func NewSearchProductsService(ctx context.Context) *SearchProductsService {
	return &SearchProductsService{ctx: ctx}
}

// Run create note info
func (s *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	// 检查搜索参数
	if req.Query == "" {
		return nil, apiErr.ProductNameRequiredErr
	}
	if req.Page <= 0 {
		return nil, apiErr.PageRequiredErr
	}
	if req.PageSize <= 0 {
		return nil, apiErr.PageSizeRequiredErr
	}

	// 调用 Meilisearch 搜索接口
	searchResp, err := meili.Client.Index("product").SearchWithContext(s.ctx, req.Query,
		&meilisearch.SearchRequest{
			Limit:  req.PageSize,
			Offset: int64(req.Page-1) * req.PageSize,
		})
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}

	// 处理搜索结果
	result := make([]*product.Product, 0)
	for _, hit := range searchResp.Hits {
		// 将搜索结果转换为 map
		doc, ok := hit.(map[string]any)
		if !ok {
			// 处理断言失败的情况，记录日志并返回错误
			klog.CtxErrorf(s.ctx, "搜索结果类型转换失败, hit: %v", hit)
			return nil, apiErr.ServiceErr
		}

		// 转换各个字段
		id, err := convertToUint32(doc["id"])
		if err != nil {
			klog.CtxErrorf(s.ctx, "ID转换失败: %v", err)
			return nil, apiErr.ServiceErr
		}
		price, err := convertToFloat32(doc["price"])
		if err != nil {
			klog.CtxErrorf(s.ctx, "价格转换失败: %v", err)
			return nil, apiErr.ServiceErr
		}
		categories, err := convertToStringSlice(doc["categories"])
		if err != nil {
			klog.CtxErrorf(s.ctx, "分类转换失败: %v", err)
			return nil, apiErr.ServiceErr
		}
		name, ok := doc["name"].(string)
		if !ok {
			klog.CtxErrorf(s.ctx, "商品名称类型错误")
			return nil, apiErr.ServiceErr
		}
		picture, ok := doc["picture"].(string)
		if !ok {
			klog.CtxErrorf(s.ctx, "商品图片类型错误")
			return nil, apiErr.ServiceErr
		}

		// 构建商品对象
		result = append(result, &product.Product{
			Id:         id,
			Name:       name,
			Price:      price,
			Picture:    picture,
			Categories: categories,
		})
	}

	// 计算总页数并返回结果
	totalPages := int64(math.Ceil(float64(searchResp.EstimatedTotalHits) / float64(req.PageSize)))
	return &product.SearchProductsResp{
		Results: result,
		Num:     totalPages,
	}, nil
}

func convertToUint32(value any) (uint32, error) {
	if floatVal, ok := value.(float64); ok {
		return uint32(floatVal), nil // 显式转换为 uint32
	}
	return 0, fmt.Errorf("value is not a float64: %T", value)
}

func convertToFloat32(value any) (float32, error) {
	if floatVal, ok := value.(float64); ok {
		return float32(floatVal), nil
	}
	return 0, fmt.Errorf("value is not a valid float64: %T", value)
}

func convertToStringSlice(value any) ([]string, error) {
	interfaceSlice, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("value is not a []interface{}: %T", value)
	}

	stringSlice := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element is not a string: %T", v)
		}
		stringSlice[i] = str
	}
	return stringSlice, nil
}
