package logic

import (
	"context"
	xerrors "github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/utils/shopping"
	"jijizhazha1024/go-mall/services/order/order"

	"jijizhazha1024/go-mall/apis/order/internal/svc"
	"jijizhazha1024/go-mall/apis/order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderReq) (resp *types.OrderDetailResp, err error) {
	userID, ok := l.ctx.Value(biz.UserIDKey).(uint32)
	if !ok {
		return nil, xerrors.New(code.AuthBlank, code.AuthBlankMsg)
	}
	res, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderRequest{
		OrderId: req.OrderID,
		UserId:  userID,
	})
	if err != nil {
		l.Logger.Errorw("call rpc GetOrder failed", logx.Field("err", err))
		return nil, xerrors.New(code.ServerError, code.ServerErrorMsg)
	}
	if res.StatusCode != code.Success {
		return nil, xerrors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = &types.OrderDetailResp{
		Order:   convertOrder2Resp(res.Order),
		Address: convertOrderAddress2Resp(res.Address),
		Items:   convertOrderItems2Resp(res.Items),
	}
	return
}

func convertOrderItems2Resp(items []*order.OrderItem) []types.OrderItemResp {
	var itemsResp []types.OrderItemResp
	for _, item := range items {
		itemsResp = append(itemsResp, types.OrderItemResp{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			ProductDesc: item.ProductDesc,
			UnitPrice:   shopping.FenToYuan(item.UnitPrice),
			Quantity:    item.Quantity,
			ItemID:      item.ItemId,
		})
	}
	return itemsResp
}

func convertOrderAddress2Resp(address *order.OrderAddress) types.OrderAddressResp {
	return types.OrderAddressResp{
		AddressID:       address.AddressId,
		RecipientName:   address.RecipientName,
		PhoneNumber:     address.PhoneNumber,
		Province:        address.Province,
		City:            address.City,
		DetailedAddress: address.DetailedAddress,
		CreatedAt:       address.CreatedAt,
		UpdatedAt:       address.UpdatedAt,
		OrderID:         address.OrderId,
	}
}

func convertOrder2Resp(o *order.Order) types.OrderResp {
	return types.OrderResp{
		OrderID:        o.OrderId,
		PreOrderID:     o.PreOrderId,
		UserID:         o.UserId,
		PaymentMethod:  int32(o.PaymentMethod),
		TransactionID:  o.TransactionId,
		PaidAt:         o.PaidAt,
		OriginalAmount: shopping.FenToYuan(o.OriginalAmount),
		DiscountAmount: shopping.FenToYuan(o.DiscountAmount),
		PayableAmount:  shopping.FenToYuan(o.PayableAmount),
		PaidAmount:     shopping.FenToYuan(o.PaidAmount),
		OrderStatus:    int32(o.OrderStatus),
		PaymentStatus:  int32(o.PaymentStatus),
		Reason:         o.Reason,
		ExpireTime:     o.ExpireTime,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}
}
