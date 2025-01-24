package logic

import (
	"context"
	"errors"
	"jijizhazha1024/go-mall/services/carts/internal/model"

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
	// 1. 查询购物车记录是否存在
	var shopCart model.Cart

	// 查找对应的购物车记录
	if result := l.svcCtx.DB.Where("product_id = ? AND user_id = ?", in.ProductId, in.UserId).First(&shopCart); result.RowsAffected == 0 {
		// 如果没有找到，返回一个“未找到”错误
		return nil, errors.New("购物车记录不存在")
	}

	// 2. 删除购物车记录
	if result := l.svcCtx.DB.Where("product_id = ? AND user_id = ?", in.ProductId, in.UserId).Delete(&model.Cart{}); result.RowsAffected == 0 {
		// 如果删除失败，返回错误
		return nil, errors.New("删除购物车记录失败")
	}

	// 3. 返回空响应
	return &carts.EmptyCartResponse{}, nil
}
