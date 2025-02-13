package middleware

import (
	"context"
	"fmt"
	"log"

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
		log.Printf("Casbin创建gorm适配器失败: %v", err)
		return nil, err
	}
	// 加载模型
	// enforcer, err := casbin.NewEnforcer("../model/rbac.conf", adapter)
	enforcer, err := casbin.NewEnforcer("rbac.conf", adapter)
	if err != nil {
		log.Printf("创建Casbin模型失败: %v", err)
		return nil, err
	}
	// 从数据库加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Printf("加载Casbin策略失败: %v", err)
		return nil, err
	}

	if err := initDefaultPolicies(enforcer); err != nil {
		log.Printf("初始化默认策略失败: %v", err)
		return nil, err
	}

	return &CasbinMiddleware{enforcer: enforcer}, nil
}

func (cm *CasbinMiddleware) Middleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从上下文中获取角色
		roleVal, exists := c.Get("role")
		var role string
		switch v := roleVal.(type) {
		case int:
			switch v {
			case 1:
				role = "admin"
			case 2:
				role = "user"
			default:
				role = "public"
			}
		case string:
			role = v
		default:
			c.AbortWithStatus(403)
			return
		}
		if !exists {
			c.AbortWithStatus(401)
			return
		}

		// 获取请求信息
		obj := string(c.Request.URI().Path())
		act := string(c.Request.Method())

		// 权限验证
		ok, err := cm.enforcer.Enforce(fmt.Sprint(role), obj, act)
		if err != nil {
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
	if _, err := enforcer.AddPolicy("admin", "*", "*"); err != nil {
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
		{"user", "/orders/*", "PUT"},
		{"user", "/orders/*", "DELETE"},
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
