package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

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

	res := &coupons.CalculateCouponResp{}

	amount := int64(0)
	coupon, err := l.svcCtx.CouponsModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res.StatusCode = code.CouponsNotExist
			res.StatusMsg = code.CouponsNotExistMsg
			return res, errors.New(code.ServerErrorMsg)
		}
		logx.Errorw("query coupons by id error", logx.Field("err", err), logx.Field("order_id", in.OrderId))
		return nil, err
	}
	discountAmount := decimal.NewFromFloat(coupon.Value)
	originAmount := decimal.NewFromInt(amount)
	minAmount := decimal.NewFromInt(int64(coupon.MinAmount))
	switch coupons.CouponType(coupon.Type) {
	// 满减
	case coupons.CouponType_COUPON_TYPE_FULL_REDUCTION:
		if originAmount.GreaterThanOrEqual(minAmount) {
			res.DiscountAmount = discountAmount.String()
			res.FinalAmount = originAmount.Sub(discountAmount).String()
		} else {
			// 不满足
		}
		//
	// 打折
	case coupons.CouponType_COUPON_TYPE_DISCOUNT:
		res.DiscountAmount = originAmount.Mul(discountAmount).Div(decimal.NewFromInt(100)).String()
	// 立减券
	case coupons.CouponType_COUPON_TYPE_FIXED_AMOUNT:
		res.DiscountAmount = discountAmount.String()

	}
	res.FinalAmount = originAmount.Sub(discountAmount).String()

	return res, nil
}
