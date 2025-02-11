package middleware

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

// NewCasbinEnforcer 创建并初始化 Casbin Enforcer
func NewCasbinEnforcer(db *gorm.DB) *CasbinMiddleware {
	// 创建 GORM 适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		hlog.Fatalf("Casbin创建gorm适配器失败: %v", err)
	}

	// 加载模型
	enforcer, err := casbin.NewEnforcer("../model/rbac.conf", adapter)
	if err != nil {
		hlog.Fatalf("创建Casbin 模型失败: %v", err)
	}

	// 从数据库加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		hlog.Fatalf("加载Casbin 策略失败: %v", err)
	}
	initDefaultPolicies(enforcer)

	return &CasbinMiddleware{enforcer: enforcer}
}

// 初始化Casbin中间件
func InitCasbinMiddleware(enforcer *casbin.Enforcer) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从上下文中获取角色
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(401, "401")
			return
		}
		// 获取请求的路径和方法
		path := string(c.Request.URI().Path())
		method := string(c.Request.Method())
		// 使用Casbin检查权限
		ok, err := enforcer.Enforce(fmt.Sprint(role), path, method)
		if err != nil {
			c.AbortWithStatusJSON(500, "500")
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, "403")
			return
		}
		c.Next(ctx)
	}
}

func (cm *CasbinMiddleware) Middleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从上下文中获取角色
		role, exists := c.Get("role")
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

func initDefaultPolicies(enforcer *casbin.Enforcer) {
	// 管理员权限
	_, _ = enforcer.AddPolicy("admin", "*", "*")

	// 公共访问权限
	_, _ = enforcer.AddPolicy("public", "/auth/token", "POST")
	_, _ = enforcer.AddPolicy("public", "/auth/verify", "POST")
	_, _ = enforcer.AddPolicy("public", "/auth/renew", "POST")
	// 用户权限
	_, _ = enforcer.AddPolicy("public", "/user/register", "POST")
	_, _ = enforcer.AddPolicy("public", "/user/login", "POST")
	_, _ = enforcer.AddPolicy("user", "/user/logout", "POST")
	_, _ = enforcer.AddPolicy("user", "/user/update", "PUT")
	_, _ = enforcer.AddPolicy("user", "/user/info", "GET")

	_, _ = enforcer.AddPolicy("admin", "/user/delete", "DELETE")
	_, _ = enforcer.AddPolicy("admin", "/user/update_role", "PUT")
	// 商品服务
	_, _ = enforcer.AddPolicy("public", "/products", "GET")
	_, _ = enforcer.AddPolicy("public", "/product", "GET")
	_, _ = enforcer.AddPolicy("public", "/product/search", "GET")

	_, _ = enforcer.AddPolicy("admin", "/product", "POST")
	_, _ = enforcer.AddPolicy("admin", "/product/upload", "POST")
	_, _ = enforcer.AddPolicy("admin", "/product", "PUT")
	_, _ = enforcer.AddPolicy("admin", "/product", "DELETE")
	// 订单服务
	_, _ = enforcer.AddPolicy("user", "/orders", "POST")
	_, _ = enforcer.AddPolicy("user", "/orders", "GET")
	_, _ = enforcer.AddPolicy("user", "/orders/*", "PUT")
	_, _ = enforcer.AddPolicy("user", "/orders/*", "DELETE")
	// 保存策略
	_ = enforcer.SavePolicy()
}
