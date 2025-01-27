package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/carts/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCartItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCartItemLogic {
	return &CreateCartItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCartItemLogic) CreateCartItem(in *carts.CartItemRequest) (*carts.CartInfoResponse, error) {
	// todo: add your logic here and delete this line
	//// 1. 查询购物车中是否已经有该商品
	//var shopCart model.Cart
	//
	//// 查询购物车中是否已有该商品，且关联查询 User 和 Product 数据
	//if result := l.svcCtx.DB.Preload("User").Preload("Product").Where(&model.Cart{ProductID: in.ProductId, UserID: in.UserId}).First(&shopCart); result.RowsAffected == 1 {
	//	// 如果购物车中已经存在该商品，则更新数量
	//	shopCart.Quantity += in.Quantity
	//	l.svcCtx.DB.Save(&shopCart)
	//} else {
	//	// 如果商品不在购物车中，创建新记录
	//	shopCart.UserID = in.UserId
	//	shopCart.ProductID = in.ProductId
	//	shopCart.Quantity = in.Quantity
	//	shopCart.Checked = false // 默认未选中
	//	l.svcCtx.DB.Save(&shopCart)
	//}
	//
	//// 2. 返回响应
	//return &carts.CartInfoResponse{
	//	Id:        shopCart.ID,
	//	UserId:    shopCart.UserID,
	//	ProductId: shopCart.ProductID,
	//	Quantity:  shopCart.Quantity,
	//	Checked:   shopCart.Checked,
	//}, nil
	return nil, nil
}
