package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCheckoutDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCheckoutDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCheckoutDetailLogic {
	return &GetCheckoutDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCheckoutDetail 获取结算详情
func (l *GetCheckoutDetailLogic) GetCheckoutDetail(in *checkout.CheckoutDetailReq) (*checkout.CheckoutDetailResp, error) {
	// 第一步：根据传入的 UserId 和 PreOrderId 获取结算记录
	checkoutRecord, err := l.svcCtx.CheckoutModel.FindOneByUserIdAndPreOrderId(l.ctx, in.UserId, in.PreOrderId)
	if err != nil {
			return nil, err
		}
		l.Logger.Errorw("查询结算记录失败", logx.Field("err", err), logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
		return &checkout.CheckoutDetailResp{
			StatusCode: 500,
			StatusMsg:  "查询结算记录失败",
		}, err
	}

	// 第二步：根据 PreOrderId 获取结算项（Items）
	checkoutItems, err := l.svcCtx.CheckoutModel.find(l.ctx, in.PreOrderId)
	if err != nil {
		l.Logger.Errorw("获取结算项失败", logx.Field("err", err), logx.Field("pre_order_id", in.PreOrderId))
		// 返回部分数据，不影响结算记录返回
	}

	// 将结算记录映射到 CheckoutOrder 结构体
	orderData := &checkout.CheckoutOrder{
		PreOrderId: checkoutRecord.PreOrderId,
		UserId:     checkoutRecord.UserId,
		Status:     checkoutRecord.Status,
		ExpireTime: checkoutRecord.ExpireTime,
		CreatedAt:  checkoutRecord.CreatedAt,
		UpdatedAt:  checkoutRecord.UpdatedAt,
		Items:      checkoutItems.Items, // 将结算项添加到 Items
	}

	// 第三步：构建响应结构体
	resp := &checkout.CheckoutDetailResp{
		StatusCode: 200, // 假设请求成功返回 200
		StatusMsg:  "成功获取结算详情",
		Data:       orderData, // 将 CheckoutOrder 作为 Data 返回
	}

	// 第四步：返回填充完的响应对象
	return resp, nil
}
}
