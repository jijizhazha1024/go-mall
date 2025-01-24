package logic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/services/carts/internal/model"

	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/carts/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCartItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartItemLogic {
	return &UpdateCartItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCartItemLogic) UpdateCartItem(in *carts.CartItemRequest) (*carts.EmptyCartResponse, error) {
	// todo: add your logic here and delete this line
	// 查询购物车记录是否存在
	var shopCart model.Cart

	// 查找对应的购物车记录
	if result := l.svcCtx.DB.Where("product_id = ? AND user_id = ?", in.ProductId, in.UserId).First(&shopCart); result.RowsAffected == 0 {
		// 如果没有找到，返回一个“未找到”错误
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	// 更新选中状态
	shopCart.Checked = in.Checked

	// 更新商品数量（如果数量大于0）
	if in.Quantity > 0 {
		shopCart.Quantity = in.Quantity
	}

	// 保存更新后的购物车记录
	if result := l.svcCtx.DB.Save(&shopCart); result.Error != nil {
		// 如果保存失败，返回数据库错误
		return nil, status.Errorf(codes.Internal, "更新购物车记录失败: %v", result.Error)
	}

	// 返回空响应，表示更新成功
	return &carts.EmptyCartResponse{}, nil
}
