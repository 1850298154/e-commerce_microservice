package service

import (
	"2501YTC/app/auth/biz/middlewares"
	models "2501YTC/app/auth/biz/model"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"github.com/dgrijalva/jwt-go"
	"time"
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
		UserId:      uint(req.UserId),
		AuthorityID: uint(req.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(), //生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "gomall",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, err
	}
	return &auth.DeliveryResp{
		Token:        token,
		RefreshToken: token, //待写
	}, nil
}
