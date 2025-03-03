package product

import (
	"bytes"
	"testing"

	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/common/test/assert"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestCreateProduct(t *testing.T) {
	h := server.Default()
	rpc.InitClient()
	h.POST("/product", CreateProduct)
	path := "/product"
	reqBody := product.CreateProductReq{
		Name:        "test1",
		Description: "test",
		Price:       1,
		Picture:     "",
		Categories:  []string{"test"},
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)} // todo: you can customize body
	header := ut.Header{}                                                 // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	assert.DeepEqual(t, 200, resp.StatusCode())
}

func TestUpdateProduct(t *testing.T) {
	h := server.Default()
	rpc.InitClient()
	h.PUT("/product", UpdateProduct)
	path := "/product"
	reqBody := product.UpdateProductReq{
		Id:          1,
		Name:        "test",
		Description: "test",
		Price:       1,
		Picture:     "",
		Categories:  []string{"test"},
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)} // todo: you can customize body
	header := ut.Header{}                                                 // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "PUT", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestDeleteProduct(t *testing.T) {
	h := server.Default()
	rpc.InitClient()
	h.DELETE("/product", DeleteProduct)
	path := "/product" // todo: you can customize query

	reqBody := product.DeleteProductReq{
		Id: 1,
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)} // todo: you can customize body
	header := ut.Header{}                                                 // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "DELETE", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestListProducts(t *testing.T) {
	h := server.Default()
	rpc.InitClient()
	h.GET("/products", ListProducts)
	path := "/products" // todo: you can customize query
	reqBody := product.ListProductsReq{
		Page:         1,
		PageSize:     10,
		CategoryName: "test",
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)} // todo: you can customize body
	header := ut.Header{}                                                 // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "GET", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestGetProduct(t *testing.T) {
	h := server.Default()
	rpc.InitClient()
	h.GET("/product", GetProduct)
	path := "/product" // todo: you can customize query
	reqBody := product.GetProductReq{
		Id: 1,
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)} // todo: you can customize body
	header := ut.Header{}                                                 // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "GET", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestSearchProducts(t *testing.T) {
	h := server.Default()
	rpc.InitClient()
	h.GET("/product/search", SearchProducts)
	path := "/product/search" // todo: you can customize query

	reqBody := product.SearchProductsReq{
		Query:    "test",
		Page:     1,
		PageSize: 10,
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)} // todo: you can customize body
	header := ut.Header{}                                                 // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "GET", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestUploadImage(t *testing.T) {
	h := server.Default()
	rpc.InitClient()
	h.POST("/product/upload", UploadImage)
	path := "/product/upload" // todo: you can customize query

	reqBody := product.UploadImageReq{
		ImageData: []byte{},
		FileName:  "",
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)} // todo: you can customize body
	header := ut.Header{}                                                 // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}
