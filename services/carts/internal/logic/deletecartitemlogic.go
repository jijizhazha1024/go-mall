package logic

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/carts/internal/svc"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCartItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCartItemLogic {
	return &DeleteCartItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCartItemLogic) DeleteCartItem(in *carts.CartItemRequest) (*carts.EmptyCartResponse, error) {
	// 调用数据库方法删除购物车商品
	err := l.svcCtx.CartsModel.DeleteCartItem(l.ctx, in.UserId, in.ProductId)
	if err != nil {
		// 判断是否是 "商品不存在" 的错误
		if strings.Contains(err.Error(), "not found") {
			l.Logger.Errorw("Cart item not found",
				logx.Field("user_id", in.UserId),
				logx.Field("product_id", in.ProductId))
			return &carts.EmptyCartResponse{
				StatusCode: code.CartItemNotFound, // 需要定义新的状态码
				StatusMsg:  code.CartItemNotFoundMsg,
			}, nil
		}

		// 其他错误，返回清除失败
		l.Logger.Errorw("Error deleting cart item",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId),
			logx.Field("product_id", in.ProductId))
		return &carts.EmptyCartResponse{
			StatusCode: code.CartClearFailed,
			StatusMsg:  code.CartClearFailedMsg,
		}, err
	}

	// 删除成功
	l.Logger.Infow("Cart item deleted successfully",
		logx.Field("user_id", in.UserId),
		logx.Field("product_id", in.ProductId))
	return &carts.EmptyCartResponse{
		StatusCode: code.Success,
		StatusMsg:  code.CartClearedMsg,
	}, nil
}
