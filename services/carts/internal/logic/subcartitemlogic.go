package logic

import (
	"context"
	"database/sql"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/cart"
	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/carts/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubCartItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubCartItemLogic {
	return &SubCartItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubCartItemLogic) SubCartItem(in *carts.CartItemRequest) (*carts.SubCartResponse, error) {
	// 1. 检查购物车中是否已有该商品
	id, exists, err := l.svcCtx.CartsModel.CheckCartItemExists(l.ctx, in.UserId, in.ProductId)
	if err != nil {
		l.Logger.Errorw("Error checking cart item existence",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId),
			logx.Field("product_id", in.ProductId))
		return &carts.SubCartResponse{
			StatusCode: code.Fail,
			StatusMsg:  code.FailMsg,
			Id:         0,
		}, err
	}

	// 如果商品不存在于购物车，返回错误
	if !exists {
		return &carts.SubCartResponse{
			StatusCode: code.CartProductInfoRetrievalFailed,
			StatusMsg:  code.CartProductInfoRetrievalFailedMsg,
			Id:         0,
		}, nil
	}

	// 2. 获取购物车中该商品的数量
	quantity, err := l.svcCtx.CartsModel.GetQuantityByUserIdAndProductId(l.ctx, in.UserId, in.ProductId)
	if err != nil {
		l.Logger.Errorw("GetQuantityByUserIdAndProductId failed",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId),
			logx.Field("product_id", in.ProductId))
		return &carts.SubCartResponse{
			StatusCode: code.CartProductQuantityInfoFailed,
			StatusMsg:  code.CartProductQuantityInfoFailedMsg,
			Id:         0,
		}, err
	}

	// 3. 判断数量是否可以减少
	if quantity <= 1 {
		// 如果商品数量为 1，则不能再减少
		return &carts.SubCartResponse{
			StatusCode: code.Fail,
			StatusMsg:  "商品数量不能小于 1",
			Id:         0,
		}, nil
	}

	// 4. 减少商品数量 1
	err = l.svcCtx.CartsModel.Update(l.ctx, &cart.Carts{
		Id: int64(id),
		UserId: sql.NullInt64{
			Int64: int64(in.UserId),
			Valid: true,
		},
		ProductId: sql.NullInt64{
			Int64: int64(in.ProductId),
			Valid: true,
		},
		Quantity: sql.NullInt64{
			Int64: int64(quantity) - 1,
			Valid: true,
		},
		Checked: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
	})
	if err != nil {
		l.Logger.Errorw("Failed to update cart item quantity",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId),
			logx.Field("product_id", in.ProductId))
		return &carts.SubCartResponse{
			StatusCode: code.CartSubFailed,
			StatusMsg:  code.CartSubFailedMsg,
			Id:         0,
		}, err
	}

	// 5. 返回成功响应
	l.Logger.Infow("Cart item quantity updated successfully",
		logx.Field("user_id", in.UserId),
		logx.Field("product_id", in.ProductId),
		logx.Field("quantity", quantity-1))

	return &carts.SubCartResponse{
		StatusCode: code.Success,
		StatusMsg:  code.CartSubMsg,
		Id:         id,
	}, nil
}
