package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"2501YTC/app/gateway/biz/dal"
	"2501YTC/app/gateway/infra/rpc"

	"github.com/joho/godotenv"

	"2501YTC/app/gateway/hertz_gen/gateway/user"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	_ = godotenv.Load("../../../.env")
	rpc.InitClient()
	dal.Init()

	h := server.Default()
	h.POST("/user/register", Register)

	testCases := []struct {
		name       string
		reqBody    user.RegisterReq
		wantStatus int
	}{
		{
			name: "valid registration",
			reqBody: user.RegisterReq{
				Email:           "user@example.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid registration - password mismatch",
			reqBody: user.RegisterReq{
				Email:           "test@user.com",
				Password:        "123456",
				ConfirmPassword: "654321",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid registration - empty email",
			reqBody: user.RegisterReq{
				Email:           "",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid registration - is banned",
			reqBody: user.RegisterReq{
				Email:           "user@example.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			body := &ut.Body{
				Body: bytes.NewBuffer(jsonBody),
				Len:  len(jsonBody),
			}
			header := ut.Header{
				Key:   "Content-Type",
				Value: "application/json",
			}

			w := ut.PerformRequest(h.Engine, "POST", "/user/register", body, header)
			resp := w.Result()

			assert.Equal(t, tc.wantStatus, resp.StatusCode())
			t.Logf("Response body: %s", string(resp.Body()))
		})
	}
}

func TestLogin(t *testing.T) {
	h := server.Default()
	h.POST("/user/login", Login)

	testCases := []struct {
		name       string
		reqBody    user.LoginReq
		wantStatus int
	}{
		{
			name: "valid login",
			reqBody: user.LoginReq{
				Email:    "test@user.com",
				Password: "123456",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid login - wrong password",
			reqBody: user.LoginReq{
				Email:    "test@user.com",
				Password: "wrongpassword",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid login - empty email",
			reqBody: user.LoginReq{
				Email:    "",
				Password: "123456",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			body := &ut.Body{
				Body: bytes.NewBuffer(jsonBody),
				Len:  len(jsonBody),
			}
			header := ut.Header{
				Key:   "Content-Type",
				Value: "application/json",
			}

			w := ut.PerformRequest(h.Engine, "POST", "/user/login", body, header)
			resp := w.Result()

			assert.Equal(t, tc.wantStatus, resp.StatusCode())
			t.Logf("Response body: %s", string(resp.Body()))
		})
	}
}

func TestLogout(t *testing.T) {
	h := server.Default()
	h.POST("/user/logout", Logout)

	testCases := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "valid logout",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsIlJvbGUiOjEsImV4cCI6MTczODkwMzE4MiwianRpIjoiMGQzMDgyNGQtNDQwNi00YmVhLWIzMTctZTA0YzJkZWVlYmYyIiwiaWF0IjoxNzM4ODk5NTgyLCJpc3MiOiJnb21hbGwifQ.ioaJpl6kYGmAohABRTeJuVs0zcf_Lj9_m2juYlGZOQY",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid logout - empty token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid logout - invalid token",
			token:      "invalidtoken",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := &ut.Body{
				Body: bytes.NewBufferString(""),
				Len:  1,
			}
			header := ut.Header{
				Key:   "Authorization",
				Value: "Bearer " + tc.token,
			}

			w := ut.PerformRequest(h.Engine, "POST", "/user/logout", body, header)
			resp := w.Result()

			assert.Equal(t, tc.wantStatus, resp.StatusCode())
			t.Logf("Response body: %s", string(resp.Body()))
		})
	}
}

func TestDeleteUser(t *testing.T) {
	h := server.Default()
	h.DELETE("/user/delete", DeleteUser)

	testCases := []struct {
		name       string
		userID     string
		token      string
		wantStatus int
	}{
		{
			name:       "valid delete",
			userID:     "3",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsIlJvbGUiOjEsImV4cCI6MTczODkwMzE4MiwianRpIjoiMGQzMDgyNGQtNDQwNi00YmVhLWIzMTctZTA0YzJkZWVlYmYyIiwiaWF0IjoxNzM4ODk5NTgyLCJpc3MiOiJnb21hbGwifQ.ioaJpl6kYGmAohABRTeJuVs0zcf_Lj9_m2juYlGZOQY",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid delete - empty user ID",
			userID:     "",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsIlJvbGUiOjEsImV4cCI6MTczODkwMzE4MiwianRpIjoiMGQzMDgyNGQtNDQwNi00YmVhLWIzMTctZTA0YzJkZWVlYmYyIiwiaWF0IjoxNzM4ODk5NTgyLCJpc3MiOiJnb21hbGwifQ.ioaJpl6kYGmAohABRTeJuVs0zcf_Lj9_m2juYlGZOQY",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid delete - empty token",
			userID:     "3",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := "/user/delete?user_id=" + tc.userID
			body := &ut.Body{
				Body: bytes.NewBufferString(""),
				Len:  1,
			}
			header := ut.Header{
				Key:   "Authorization",
				Value: "Bearer " + tc.token,
			}

			w := ut.PerformRequest(h.Engine, "DELETE", path, body, header)
			resp := w.Result()

			assert.Equal(t, tc.wantStatus, resp.StatusCode())
			t.Logf("Response body: %s", string(resp.Body()))
		})
	}
}

func TestUpdateUser(t *testing.T) {
	h := server.Default()
	h.PUT("/user/update", UpdateUser)

	testCases := []struct {
		name       string
		reqBody    map[string]any
		wantStatus int
	}{
		{
			name: "valid update",
			reqBody: map[string]any{
				"name": "new_name",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid update - empty body",
			reqBody:    map[string]any{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			body := &ut.Body{
				Body: bytes.NewBuffer(jsonBody),
				Len:  len(jsonBody),
			}
			header := ut.Header{
				Key:   "Content-Type",
				Value: "application/json",
			}

			w := ut.PerformRequest(h.Engine, "PUT", "/user/update", body, header)
			resp := w.Result()

			assert.Equal(t, tc.wantStatus, resp.StatusCode())
			t.Logf("Response body: %s", string(resp.Body()))
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	h := server.Default()
	h.GET("/user/info", GetUserInfo)

	testCases := []struct {
		name       string
		wantStatus int
	}{
		{
			name:       "valid get user info",
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := &ut.Body{
				Body: bytes.NewBufferString(""),
				Len:  1,
			}
			header := ut.Header{}

			w := ut.PerformRequest(h.Engine, "GET", "/user/info", body, header)
			resp := w.Result()

			assert.Equal(t, tc.wantStatus, resp.StatusCode())
			t.Logf("Response body: %s", string(resp.Body()))
		})
	}
}

func TestUpdateUserRole(t *testing.T) {
	h := server.Default()
	h.PUT("/user/update_role", UpdateUserRole)

	testCases := []struct {
		name       string
		reqBody    map[string]any
		wantStatus int
	}{
		{
			name: "valid update role",
			reqBody: map[string]any{
				"user_id": uint32(1),
				"role":    2,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid update role - wrong user_id",
			reqBody: map[string]any{
				"user_id": uint32(123),
				"role":    2,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			body := &ut.Body{
				Body: bytes.NewBuffer(jsonBody),
				Len:  len(jsonBody),
			}
			header := ut.Header{
				Key:   "Content-Type",
				Value: "application/json",
			}

			w := ut.PerformRequest(h.Engine, "PUT", "/user/update_role", body, header)
			resp := w.Result()

			assert.Equal(t, tc.wantStatus, resp.StatusCode())
			t.Logf("Response body: %s", string(resp.Body()))
		})
	}
}
