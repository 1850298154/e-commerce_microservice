package service

import (
	"context"
	"fmt"
	"strconv"

	"2501YTC/app/checkout/infra/rpc"
	"2501YTC/rpc_gen/kitex_gen/cart"
	"2501YTC/rpc_gen/kitex_gen/checkout"
	"2501YTC/rpc_gen/kitex_gen/order"
	"2501YTC/rpc_gen/kitex_gen/payment"
	"2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

// Run create note info

/*
	Run

// 1. get cart
// 2. calculate cart
// 3. create order
// 4. empty cart
// 5. pay
// 6. change order result
// 7. finish
*/
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	fmt.Println("CheckoutService.Run")
	// Finish your business logic.

	// 1. get cart
	// Idempotent
	// get cart
	cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
	// cartResult.Cart.Items

	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005001, err.Error())
	}
	if cartResult == nil || cartResult.Cart.Items == nil {
		return nil, kerrors.NewGRPCBizStatusError(5004001, "cart is empty")
	}

	// 2. calculate cart
	var (
		oi    []*order.OrderItem
		total float32
	)
	for _, cartItem := range cartResult.Cart.Items {
		productResp, resultErr := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
			Id: cartItem.ProductId,
		})

		if resultErr != nil {
			return nil, resultErr
		}

		if productResp.Product == nil {
			continue
		}

		p := productResp.Product.Price

		cost := p * float32(cartItem.Quantity)
		total += cost
		oi = append(oi, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId: cartItem.ProductId,
				Quantity:  cartItem.Quantity,
			},
			Cost: cost,
		})
	}
	fmt.Println("total", total)

	// 3. create order
	orderReq := &order.PlaceOrderReq{
		UserId:       req.UserId,
		UserCurrency: "USD",
		OrderItems:   oi,
		Email:        req.Email,
	}
	if req.Address != nil {
		addr := req.Address
		zipCodeInt, _ := strconv.Atoi(addr.ZipCode)
		orderReq.Address = &order.Address{
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			Country:       addr.Country,
			State:         addr.State,
			ZipCode:       int32(zipCodeInt),
		}
	}
	orderResult, err := rpc.OrderClient.PlaceOrder(s.ctx, orderReq)
	if err != nil {
		err = fmt.Errorf("PlaceOrder.err:%v", err)
		return
	}
	klog.Info("orderResult", orderResult)
	fmt.Println("orderResult", orderResult)

	// 4. empty cart
	emptyResult, err := rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
	if err != nil {
		err = fmt.Errorf("EmptyCart.err:%v", err)
		return
	}
	klog.Info("emptyResult")
	klog.Info(emptyResult)
	fmt.Println("emptyResult")
	fmt.Println(emptyResult)

	// 5. pay
	// ==charge
	var orderId string
	if orderResult != nil || orderResult.Order != nil {
		orderId = orderResult.Order.OrderId
	}

	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
		},
	}

	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
	if err != nil {
		return nil, err
	}
	klog.Info("paymentResult")
	klog.Info(paymentResult)
	fmt.Println("paymentResult")
	fmt.Println(paymentResult)

	// 6. change order result
	// change order state
	_, err = rpc.OrderClient.MarkOrderPaid(s.ctx, &order.MarkOrderPaidReq{UserId: req.UserId, OrderId: orderId})
	if err != nil {
		klog.Error(err)
		return
	}

	// 7. finish
	resp = &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
	}
	fmt.Println("normal return ......")
	return
}
