package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/cart"
	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/carts/internal/svc"
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
	// 定义响应对象
	var rsp carts.CartItemListResponse

	// 获取购物车模型
	cartModel := l.svcCtx.Mysql
	// 查询当前用户的购物车数据
	var shopCarts []cart.Carts
	query := `
	SELECT id, created_at,updated_at,deleted_at,user_id, product_id, quantity, checked
   	FROM carts
   	WHERE user_id = ? AND deleted_at IS NULL
    `
	err := cartModel.QueryRows(&shopCarts, query, in.Id)
	if err != nil {
		logx.Errorf("Error occurred while querying carts: %v", err)
		if err == sqlx.ErrNotFound {
			return &rsp, nil // 没有找到数据，返回空响应
		}
		return nil, err // 其他错误，返回错误
	}
	// 设置响应中的总数
	rsp.Total = int32(len(shopCarts))

	// 构建响应数据
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &carts.CartInfoResponse{
			Id:        int32(shopCart.Id),
			UserId:    int32(shopCart.UserId.Int64),
			ProductId: int32(shopCart.ProductId.Int64),
			Quantity:  int32(shopCart.Quantity.Int64),
			Checked:   shopCart.Checked.Int64 == 1,
		})
	}

	return &rsp, nil
}
