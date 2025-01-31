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
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			Issuer:    "gomall",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, err
	}

	refreshClaims := claims
	refreshClaims.StandardClaims.ExpiresAt *= time.Now().Add(7 * 24 * time.Hour).Unix() // refreshtoken 7天
	refreshToken, err := j.CreateToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	// 数据库保存
	tokenRecord := models.Token{
		UserID:         uint(req.UserId),
		Role:           req.Role,
		Token:          token,
		RefreshToken:   refreshToken,
		AccessExpires:  time.Now().Add(1 * time.Hour).Unix(),      // 30天过期
		RefreshExpires: time.Now().Add(7 * 24 * time.Hour).Unix(), // 60天过期
	}
	tokenQuery := models.NewTokenQuery(s.ctx, mysql.DB)
	_, err = tokenQuery.Create(tokenRecord)
	if err != nil {
		return nil, err
	}

	return &auth.DeliveryResp{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
