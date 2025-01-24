package service

import (
	"2501YTC/app/auth/biz/middlewares"
	models "2501YTC/app/auth/biz/model"
	"2501YTC/app/order/biz/dal/mysql"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type DeliverTokenByRPCService struct {
	ctx context.Context
} // NewDeliverTokenByRPCService new DeliverTokenByRPCService

func NewDeliverTokenByRPCService(ctx context.Context) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		UserId: uint(req.UserId),
		Role:   req.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(), // 生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "gomall",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, err
	}

	claims2 := claims
	claims2.StandardClaims.NotBefore *= 2 // refreshtoken 增加的时间
	refresh_token, err := j.CreateToken(claims2)
	if err != nil {
		return nil, err
	}

	// 数据库保存
	tokenRecord := models.Token{
		UserID:         uint(req.UserId),
		Role:           req.Role,
		Token:          token,
		RefreshToken:   refresh_token,
		AccessExpires:  time.Now().Unix() + 60*60*24*30, // 30天过期
		RefreshExpires: time.Now().Unix() + 60*60*24*60, // 60天过期
	}
	tokenQuery := models.NewTokenQuery(s.ctx, mysql.DB)
	_, err = tokenQuery.Create(tokenRecord)
	if err != nil {
		return nil, err
	}

	return &auth.DeliveryResp{
		Token:        token,
		RefreshToken: refresh_token,
	}, nil
}
