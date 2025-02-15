package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/carts/carts"

	"jijizhazha1024/go-mall/apis/carts/internal/svc"
	"jijizhazha1024/go-mall/apis/carts/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCartItemLogic {
	return &CreateCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCartItemLogic) CreateCartItem(req *types.CreateCartReq) (resp *types.CreateCartResp, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.CartRpc.CreateCartItem(l.ctx, &carts.CartItemRequest{
		UserId:       req.UserId,
		ProductId:    req.ProductId,
		ProductName:  req.ProductName,
		ProductImage: req.ProuctImage,
		ProductPrice: req.ProductPrice,
		Quantity:     req.Quantity,
		Checked:      req.Checked,
	})
	if err != nil {
		if res != nil && res.StatusCode != code.Success {
			// 处理用户级别info 错误
			return nil, errors.New(int(res.StatusCode), res.StatusMsg)
		}
		l.Logger.Errorf("call rpc CreateCartItem failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}
	resp = &types.CreateCartResp{
		Id: res.Id,
	}
	return
}
