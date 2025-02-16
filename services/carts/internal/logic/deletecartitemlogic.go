package logic

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/carts/internal/svc"

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
	// todo: add your logic here and delete this line
	err := l.svcCtx.CartsModel.DeleteCartItem(l.ctx, in.UserId, in.ProductId)
	if err != nil {
		l.Logger.Errorw("Error deleting cart item",
			logx.Field("err", err),
			logx.Field("user_id", in.Id),
			logx.Field("product_id", in.ProductId))
		return &carts.EmptyCartResponse{
			StatusCode: code.CartClearFailed,
			StatusMsg:  code.CartClearFailedMsg,
		}, err
	} else {
		return &carts.EmptyCartResponse{
			StatusCode: code.Success,
			StatusMsg:  code.CartClearedMsg,
		}, nil
	}
}
