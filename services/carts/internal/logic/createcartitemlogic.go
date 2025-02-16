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
	// todo: add your logic here and delete this line\

	id, exists, err := l.svcCtx.CartsModel.CheckCartItemExists(l.ctx, in.UserId, in.ProductId)
	if err != nil {
		l.Logger.Errorw("Error checking cart item existence",
			logx.Field("err", err),
			logx.Field("user_id", in.Id),
			logx.Field("product_id", in.ProductId))
		return &carts.CreateCartResponse{
			StatusCode: code.Fail,
			StatusMsg:  code.FailMsg,
			Id:         0,
		}, err
	} else if exists {
		quantity, err := l.svcCtx.CartsModel.GetQuantityByUserIdAndProductId(l.ctx, in.UserId, in.ProductId)
		if err != nil {
			logx.Errorf("GetQuantityByUserIdAndProductId err: %v", err)
		}
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
				Int64: int64(quantity) + 1,
				Valid: true,
			},
			Checked: sql.NullInt64{
				Int64: 1,
				Valid: true,
			},
		})
		if err != nil {
			l.Logger.Errorw("Shopcart create item err",
				logx.Field("err", err),
				logx.Field("user_id", in.Id),
				logx.Field("product_id", in.ProductId))
			return &carts.CreateCartResponse{
				StatusCode: code.CartCreationFailed,
				StatusMsg:  code.CartCreationFailedMsg,
				Id:         0,
			}, err
		}
		// Return success response
		return &carts.CreateCartResponse{
			StatusCode: code.Success,
			StatusMsg:  code.CartCreatedMsg,
			Id:         id,
		}, nil
	}

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
			Int64: int64(in.Quantity) + 1,
			Valid: true,
		},
		Checked: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
	})
	if err != nil {
		l.Logger.Errorw("Shopcart create item err",
			logx.Field("err", err),
			logx.Field("user_id", in.Id),
			logx.Field("product_id", in.ProductId))
		return &carts.CreateCartResponse{
			StatusCode: code.CartCreationFailed,
			StatusMsg:  code.CartCreationFailedMsg,
			Id:         0,
		}, err
	}
	rowsAffected, err := result.RowsAffected()
	lastInsertId, _ := result.LastInsertId()
	if rowsAffected == 0 {
		return &carts.CreateCartResponse{
			StatusCode: code.CartCreationFailed,
			StatusMsg:  code.CartCreationFailedMsg,
			Id:         int32(lastInsertId),
		}, err
	}

	return &carts.CreateCartResponse{
		StatusCode: code.Success,
		StatusMsg:  code.CartCreatedMsg,
		Id:         int32(lastInsertId),
	}, nil
}
