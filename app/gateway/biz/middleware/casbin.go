package middleware

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

// NewCasbinEnforcer 创建并初始化 Casbin Enforcer
func NewCasbinEnforcer(db *gorm.DB) (*CasbinMiddleware, error) {
	// 创建 GORM 适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		hlog.Error("Casbin创建gorm适配器失败: %v", err)
		return nil, err
	}
	// 加载模型
	fmt.Println(os.Getwd())
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(filename))
	fmt.Println(filepath.Join(basePath, "model", "rbac.conf"))
	var modelPath string
	if env := GetEnv(); env != "online" {
		modelPath = filepath.Join(basePath, "model", "rbac.conf")
	} else {
		modelPath = "rbac.conf"
	}
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		hlog.Error("创建Casbin模型失败: %v", err)
		return nil, err
	}

	enforcer.AddFunction("RegexMatch", RegexMatch)
	// 从数据库加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		hlog.Error("加载Casbin策略失败: %v", err)
		return nil, err
	}

	if err := initDefaultPolicies(enforcer); err != nil {
		hlog.Error("初始化默认策略失败: %v", err)
		return nil, err
	}

	return &CasbinMiddleware{enforcer: enforcer}, nil
}

func (cm *CasbinMiddleware) Middleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var role string
		fmt.Println(c.Get("user_id"))
		// 从上下文中获取角色
		roleVal, exists := c.Get("role")
		if !exists {
			role = "public" // 如果没有角色，则默认为 public
		} else {
			switch v := roleVal.(type) {
			case uint32:
				if v == 1 {
					role = "admin"
				} else if v == 2 {
					role = "user"
				} else {
					role = "public"
				}
			case string:
				role = v
			default:
				role = "public" // 默认设置为 public
			}
		}
		// 获取请求信息
		obj := string(c.Request.URI().Path())
		act := string(c.Request.Method())

		// 权限验证
		ok, err := cm.enforcer.Enforce(fmt.Sprint(role), obj, act)
		if err != nil {
			hlog.Error("Casbin 权限验证失败: %v", err)
			c.AbortWithStatus(500)
			return
		}

		if !ok {
			c.AbortWithStatus(403)
			return
		}
		c.Next(ctx)
	}
}

func initDefaultPolicies(enforcer *casbin.Enforcer) error {
	// 管理员权限
	if _, err := enforcer.AddPolicy("admin", ".*", ".*"); err != nil {
		return fmt.Errorf("添加管理员策略失败: %w", err)
	}

	// 公共访问权限
	policies := [][]string{
		{"public", "/auth/token", "POST"},
		{"public", "/auth/verify", "POST"},
		{"public", "/auth/renew", "POST"},
		{"public", "/user/register", "POST"},
		{"public", "/user/login", "POST"},
		{"user", "/user/logout", "POST"},
		{"user", "/user/update", "PUT"},
		{"user", "/user/info", "GET"},
		{"admin", "/user/delete", "DELETE"},
		{"admin", "/user/update_role", "PUT"},
		{"public", "/products", "GET"},
		{"public", "/product", "GET"},
		{"public", "/product/search", "GET"},
		{"admin", "/product", "POST"},
		{"admin", "/product/upload", "POST"},
		{"admin", "/product", "PUT"},
		{"admin", "/product", "DELETE"},
		{"user", "/orders", "POST"},
		{"user", "/orders", "GET"},
		{"user", "/orders/.*", "PUT"},
		{"user", "/orders/.*", "DELETE"},
		{"public", "/checkout", "POST"},
		{"public", "/checkout", "GET"},
		{"public", "/checkout/.*", "PUT"},
		{"public", "/checkout/.*", "DELETE"},
		{"user", "/ai/.*", "POST"},
		{"admin", "/ai/.*", "POST"},
	}

	for _, p := range policies {
		if _, err := enforcer.AddPolicy(p[0], p[1], p[2]); err != nil {
			return fmt.Errorf("添加策略%v失败:%w", p, err)
		}
	}
	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存Casbin策略失败:%w", err)
	}
	return nil
}

func RegexMatch(args ...any) (any, error) {
	key, ok := args[0].(string)
	if !ok {
		err := errors.New("key错误")
		hlog.Error(err)
		return nil, err
	}

	pattern, ok := args[1].(string)
	if !ok {
		err := errors.New("pattern错误")
		hlog.Error(err)
		return nil, errors.New("pattern错误")
	}
	matched, err := regexp.MatchString("^"+pattern+"$", key)
	if err != nil {
		return false, err
	}
	return matched, nil
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "dev"
	}
	return e
}
