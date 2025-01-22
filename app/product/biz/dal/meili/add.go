package meili

import (
	"strings"

	"2501YTC/app/product/biz/model"
)

func Add(products []model.Product) error {
	documents := make([]map[string]any, 0)
	// 遍历产品列表，将每个产品转换为 map 并追加到切片中
	for _, p := range products {
		document := map[string]any{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"picture":     p.Picture,
			"price":       p.Price,
			"categories":  strings.Split(p.Categories, "+"),
		}
		documents = append(documents, document)
	}

	_, err := Client.Index("product").AddDocuments(documents)
	if err != nil {
		return err
	}
	return nil
}
