package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/apis/carts/internal/svc"
	"jijizhazha1024/go-mall/apis/carts/internal/types"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/carts/carts"
)

type CartItemListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartItemListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartItemListLogic {
	return &CartItemListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartItemListLogic) CartItemList(req *types.UserInfo) (resp *types.CartItemListResp, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.CartRpc.CartItemList(l.ctx, &carts.UserInfo{
		Id: req.Id,
	})
	if err != nil {
		if res != nil && res.StatusCode != code.Success {
			// 处理用户级别info 错误
			return nil, errors.New(int(res.StatusCode), res.StatusMsg)
		}
		l.Logger.Errorf("call rpc CartItemList failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}
	resp = &types.CartItemListResp{
		Total:    res.Total,
		CartInfo: ConvertCartInfoResponse(res.Data),
	}
	return
}
