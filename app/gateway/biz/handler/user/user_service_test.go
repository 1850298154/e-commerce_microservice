package user

import (
	"2501YTC/app/gateway/hertz_gen/gateway/user"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	// "github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestRegister(t *testing.T) {
	h := server.Default()
	h.POST("/user/register", Register)

	path := "/user/register"
	// 创建请求体
	reqBody := user.RegisterReq{
		Email:           "test@user.com",
		Password:        "123456",
		ConfirmPassword: "123456",
	}

	// 将请求体序列化为 JSON
	jsonBody, _ := json.Marshal(reqBody)

	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)}
	header := ut.Header{}
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestLogin(t *testing.T) {
	h := server.Default()
	h.POST("/user/login", Login)
	path := "/user/login"
	// 创建请求体
	reqBody := user.LoginReq{
		Email:    "test@user.com",
		Password: "123456",
	}

	// 将请求体序列化为 JSON
	jsonBody, _ := json.Marshal(reqBody)

	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)}
	header := ut.Header{}
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestLogout(t *testing.T) {
	h := server.Default()
	h.POST("/user/logout", Logout)
	path := "/user/logout"
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsIlJvbGUiOjEsImV4cCI6MTczODkwMzE4MiwianRpIjoiMGQzMDgyNGQtNDQwNi00YmVhLWIzMTctZTA0YzJkZWVlYmYyIiwiaWF0IjoxNzM4ODk5NTgyLCJpc3MiOiJnb21hbGwifQ.ioaJpl6kYGmAohABRTeJuVs0zcf_Lj9_m2juYlGZOQY"
	header := ut.Header{
		"Authorization",
		"Bearer " + token,
	}
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestDeleteUser(t *testing.T) {
	h := server.Default()
	h.DELETE("/user/delete", DeleteUser)
	path := "/user/delete?user_id=3"
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsIlJvbGUiOjEsImV4cCI6MTczODkwMzE4MiwianRpIjoiMGQzMDgyNGQtNDQwNi00YmVhLWIzMTctZTA0YzJkZWVlYmYyIiwiaWF0IjoxNzM4ODk5NTgyLCJpc3MiOiJnb21hbGwifQ.ioaJpl6kYGmAohABRTeJuVs0zcf_Lj9_m2juYlGZOQY"
	header := ut.Header{
		"Authorization",
		"Bearer " + token,
	}
	w := ut.PerformRequest(h.Engine, "DELETE", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestUpdateUser(t *testing.T) {
	h := server.Default()
	h.PUT("/user/update", UpdateUser)
	path := "/user/update"                                    // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "PUT", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestGetUserInfo(t *testing.T) {
	h := server.Default()
	h.GET("/user/info", GetUserInfo)
	path := "/user/info"                                      // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "GET", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestUpdateUserRole(t *testing.T) {
	h := server.Default()
	h.PUT("/user/update_role", UpdateUserRole)
	path := "/user/update_role"                               // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "PUT", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}
