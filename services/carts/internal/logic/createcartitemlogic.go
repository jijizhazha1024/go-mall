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

func (l *CreateCartItemLogic) CreateCartItem(in *carts.CartItemRequest) (*carts.CreateCartResponse, error) {
	// 1. 检查商品是否已存在于购物车
	id, exists, err := l.svcCtx.CartsModel.CheckCartItemExists(l.ctx, in.UserId, in.ProductId)
	if err != nil {
		l.Logger.Errorw("Error checking cart item existence",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId),
			logx.Field("product_id", in.ProductId))
		return &carts.CreateCartResponse{
			StatusCode: code.Fail,
			StatusMsg:  code.FailMsg,
			Id:         0,
		}, err
	}

	// 2. 如果商品已存在于购物车，则更新商品数量
	if exists {
		// 获取当前商品数量
		quantity, err := l.svcCtx.CartsModel.GetQuantityByUserIdAndProductId(l.ctx, in.UserId, in.ProductId)
		if err != nil {
			l.Logger.Errorw("Failed to get cart item quantity",
				logx.Field("err", err),
				logx.Field("user_id", in.UserId),
				logx.Field("product_id", in.ProductId))
			return &carts.CreateCartResponse{
				StatusCode: code.CartProductQuantityInfoFailed,
				StatusMsg:  code.CartProductQuantityInfoFailedMsg,
				Id:         0,
			}, err
		}

		// 增加商品数量
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
				Int64: int64(quantity) + 1, // 增加数量
				Valid: true,
			},
			Checked: sql.NullInt64{
				Int64: 1,
				Valid: true,
			},
		})

		// 错误处理
		if err != nil {
			l.Logger.Errorw("Failed to update cart item quantity",
				logx.Field("err", err),
				logx.Field("user_id", in.UserId),
				logx.Field("product_id", in.ProductId))
			return &carts.CreateCartResponse{
				StatusCode: code.CartCreationFailed,
				StatusMsg:  code.CartCreationFailedMsg,
				Id:         0,
			}, err
		}

		// 成功返回
		l.Logger.Infow("Cart item updated successfully",
			logx.Field("user_id", in.UserId),
			logx.Field("product_id", in.ProductId),
			logx.Field("quantity", quantity+1))

		return &carts.CreateCartResponse{
			StatusCode: code.Success,
			StatusMsg:  code.CartCreatedMsg,
			Id:         id,
		}, nil
	}

	// 3. 如果商品不存在于购物车，则插入新记录
	result, err := l.svcCtx.CartsModel.Insert(l.ctx, &cart.Carts{
		UserId: sql.NullInt64{
			Int64: int64(in.UserId),
			Valid: true,
		},
		ProductId: sql.NullInt64{
			Int64: int64(in.ProductId),
			Valid: true,
		},
		Quantity: sql.NullInt64{
			Int64: int64(in.Quantity) + 1, // 初始数量为 1
			Valid: true,
		},
		Checked: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
	})

	// 错误处理
	if err != nil {
		l.Logger.Errorw("Failed to create cart item",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId),
			logx.Field("product_id", in.ProductId))
		return &carts.CreateCartResponse{
			StatusCode: code.CartCreationFailed,
			StatusMsg:  code.CartCreationFailedMsg,
			Id:         0,
		}, err
	}

	// 获取操作结果
	rowsAffected, err := result.RowsAffected()
	lastInsertId, _ := result.LastInsertId()

	if rowsAffected == 0 {
		return &carts.CreateCartResponse{
			StatusCode: code.CartCreationFailed,
			StatusMsg:  code.CartCreationFailedMsg,
			Id:         int32(lastInsertId),
		}, err
	}

	// 成功返回
	l.Logger.Infow("Cart item created successfully",
		logx.Field("user_id", in.UserId),
		logx.Field("product_id", in.ProductId),
		logx.Field("cart_id", lastInsertId))

	return &carts.CreateCartResponse{
		StatusCode: code.Success,
		StatusMsg:  code.CartCreatedMsg,
		Id:         int32(lastInsertId),
	}, nil
}
