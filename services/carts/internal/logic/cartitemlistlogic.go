package logic

import (
	"context"
	"jijizhazha1024/go-mall/services/carts/internal/model"

	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/carts/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartItemListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartItemListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartItemListLogic {
	return &CartItemListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CartItemListLogic) CartItemList(in *carts.UserInfo) (*carts.CartItemListResponse, error) {
	// todo: add your logic here and delete this line
	var rsp carts.CartItemListResponse

	// 查询当前用户的购物车数据
	var shopCarts []model.Cart
	if result := l.svcCtx.DB.Where(&model.Cart{UserID: in.Id}).Find(&shopCarts); result.Error != nil {
		// 如果数据库查询出错，返回错误
		return nil, result.Error
	}

	// 设置响应中的总数
	rsp.Total = int32(len(shopCarts))

	// 构建响应数据
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &carts.CartInfoResponse{
			Id:        shopCart.ID,
			UserId:    shopCart.UserID,
			ProductId: shopCart.ProductID,
			Quantity:  shopCart.Quantity,
			Checked:   shopCart.Checked,
		})
	}

	// 返回响应
	return &rsp, nil
}
