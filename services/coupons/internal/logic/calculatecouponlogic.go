package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"
	"jijizhazha1024/go-mall/services/product/productcatalogservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCalculateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateCouponLogic {
	return &CalculateCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CalculateCoupon 计算优惠卷折扣价格
func (l *CalculateCouponLogic) CalculateCoupon(in *coupons.CalculateCouponReq) (*coupons.CalculateCouponResp, error) {
	// 根据商品id和数量计算总结价进行优惠

	// check coupon exist，is usable
	res := &coupons.CalculateCouponResp{}
	one, err := l.svcCtx.CouponsModel.FindOne(l.ctx, in.CouponId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			// 并不属于异常，提示用户优惠券不存在，且返回信息
			res.StatusCode = uint32(code.CouponsNotExist)
			res.StatusMsg = code.CouponsNotExistMsg
			return res, nil
		}
		logx.Errorw("query coupons by id error", logx.Field("err", err))
		return nil, err
	}

	// 计算总价
	var totalPrice int64
	for _, item := range in.Items {
		product, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &productcatalogservice.GetProductReq{
			Id: uint32(item.ProductId),
		})
		if err != nil {
			logx.Errorw("call rpc ProductRpc.GetProduct failed", logx.Field("err", err))
			return nil, err
		}
		if product.StatusCode != code.Success {
			// 信息不完整，商品不存在
			res.StatusCode = uint32(code.ProductNotFoundInventory)
			res.StatusMsg = code.ProductNotFoundInventoryMsg
			return res, nil
		}
		totalPrice += product.Product.Price * int64(item.Quantity)
	}
	// 根据优惠卷类型计算优惠
	var discountAmount int64 // 分
	res.IsUsable = true
	switch coupons.CouponType(one.Type) {
	// 满减券
	case coupons.CouponType_COUPON_TYPE_FULL_REDUCTION:
		if one.MinAmount > totalPrice {
			// 优惠券不满足使用条件
			res.IsUsable = false
			res.UnusableReason = "优惠券不满足使用条件"
			return res, nil
		}
		discountAmount = one.Value
		res.CouponType = biz.CouponTypeFullReduction
	// 折扣券
	case coupons.CouponType_COUPON_TYPE_DISCOUNT:
		// 优惠券无效
		if one.Value <= 0 || one.Value >= 100 {
			res.IsUsable = false
			res.UnusableReason = "无效的折扣率"
			return res, nil
		}
		discountAmount = totalPrice * (100 - one.Value) / 100
		res.CouponType = biz.CouponTypeDiscount
	//	立减券
	case coupons.CouponType_COUPON_TYPE_FIXED_AMOUNT:
		discountAmount = one.Value
		res.CouponType = biz.CouponTypeNoThreshold
	}
	res.DiscountAmount = discountAmount
	res.FinalAmount = totalPrice - discountAmount
	res.OriginAmount = totalPrice

	// 用户是否有改优惠券
	return res, nil
}
