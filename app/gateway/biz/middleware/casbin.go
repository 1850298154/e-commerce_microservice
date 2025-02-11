package middleware

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
	"log"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

// NewCasbinEnforcer 创建并初始化 Casbin Enforcer
func NewCasbinEnforcer(db *gorm.DB) *CasbinMiddleware {
	// 创建 GORM 适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Casbin创建gorm适配器失败: %v", err)

	}

	// 加载模型
	enforcer, err := casbin.NewEnforcer("../model/rbac.conf", adapter)
	if err != nil {
		log.Fatalf("创建Casbin 模型失败: %v", err)

	}

	// 从数据库加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("加载Casbin 策略失败: %v", err)

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
	enforcer.AddPolicy("admin", "*", "*")

	// 公共访问权限
	enforcer.AddPolicy("public", "/auth/token", "POST")
	enforcer.AddPolicy("public", "/auth/verify", "POST")
	enforcer.AddPolicy("public", "/auth/renew", "POST")
	// 用户权限
	enforcer.AddPolicy("public", "/user/register", "POST")
	enforcer.AddPolicy("public", "/user/login", "POST")
	enforcer.AddPolicy("user", "/user/logout", "POST")
	enforcer.AddPolicy("user", "/user/update", "PUT")
	enforcer.AddPolicy("user", "/user/info", "GET")

	enforcer.AddPolicy("admin", "/user/delete", "DELETE")
	enforcer.AddPolicy("admin", "/user/update_role", "PUT")
	// 商品服务
	enforcer.AddPolicy("public", "/products", "GET")
	enforcer.AddPolicy("public", "/product", "GET")
	enforcer.AddPolicy("public", "/product/search", "GET")

	enforcer.AddPolicy("admin", "/product", "POST")
	enforcer.AddPolicy("admin", "/product/upload", "POST")
	enforcer.AddPolicy("admin", "/product", "PUT")
	enforcer.AddPolicy("admin", "/product", "DELETE")
	// 订单服务
	enforcer.AddPolicy("user", "/orders", "POST")
	enforcer.AddPolicy("user", "/orders", "GET")
	enforcer.AddPolicy("user", "/orders/*", "PUT")
	enforcer.AddPolicy("user", "/orders/*", "DELETE")
	// 保存策略
	enforcer.SavePolicy()
}
