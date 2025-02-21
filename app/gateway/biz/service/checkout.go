package service

import (
	"context"

	checkout "2501YTC/app/gateway/hertz_gen/gateway/checkout"
	common "2501YTC/app/gateway/hertz_gen/gateway/common"
	"2501YTC/app/gateway/infra/rpc"
	rpccheckout "2501YTC/rpc_gen/kitex_gen/checkout"
	rpcpayment "2501YTC/rpc_gen/kitex_gen/payment"

	"github.com/cloudwego/hertz/pkg/app"
)

type CheckoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCheckoutService(ctx context.Context, requestContext *app.RequestContext) *CheckoutService {
	return &CheckoutService{RequestContext: requestContext, Context: ctx}
}

func (h *CheckoutService) Run(req *checkout.CheckoutReq) (resp *common.Empty, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	userId := 1
	_, err = rpc.CheckoutClient.Checkout(h.Context, &rpccheckout.CheckoutReq{
		UserId:    uint32(userId),
		Email:     req.Email,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Address: &rpccheckout.Address{
			Country:       req.Country,
			ZipCode:       req.Zipcode,
			City:          req.City,
			State:         req.Province,
			StreetAddress: req.Street,
		},
		CreditCard: &rpcpayment.CreditCardInfo{
			CreditCardNumber:          req.CardNum,
			CreditCardExpirationYear:  req.ExpirationYear,
			CreditCardExpirationMonth: req.ExpirationMonth,
			CreditCardCvv:             req.Cvv,
		},
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
