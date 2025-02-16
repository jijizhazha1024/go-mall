package logic

import (
	"context"
	"database/sql"
	"errors"
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
	logx.Debug(coupon)
	return res, nil
}
