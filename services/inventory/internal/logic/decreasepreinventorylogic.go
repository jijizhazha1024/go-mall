package logic

import (
	"context"
	"fmt"

	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DecreasePreInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecreasePreInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreasePreInventoryLogic {
	return &DecreasePreInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DecreaseInventory 预扣减库存，此时并非真实扣除库存，而是在缓存进行--操作
func (l *DecreasePreInventoryLogic) DecreasePreInventory(in *inventory.InventoryReq) (*inventory.InventoryResp, error) {

	resp := &inventory.InventoryResp{}
	// 构建幂等锁Key（用户ID+预订单ID）
	lockKey := fmt.Sprintf("%s:%d:%s", biz.InventoryDeductLockPrefix, in.UserId, in.PreOrderId)
	//准备参数
	keys := make([]string, len(in.Items)+1)
	args := make([]interface{}, len(in.Items)+1)
	keys[0] = lockKey
	args[0] = in.PreOrderId // 构造库存Key列表

	for i, item := range in.Items {
		if item.Quantity <= 0 {
			l.Logger.Infow("商品数量不合法",
				logx.Field("product_id", item.ProductId))
			resp.StatusCode = code.InvalidParams
			resp.StatusMsg = code.InvalidParamsMsg
			return resp, nil
		}
		productKey := fmt.Sprintf("%s:%d", biz.InventoryProductKey, item.ProductId)
		keys[i+1] = productKey
		args[i+1] = item.Quantity
	}

	// 执行Lua脚本（使用go-zero的Evalsah方法）
	val, err := l.svcCtx.Rdb.EvalSha(l.svcCtx.DecreaseInventoryShal, keys, args)
	if err != nil {

		l.Logger.Errorw("LUA脚本执行失败",
			logx.Field("error", err),
			logx.Field("pre_order_id", in.PreOrderId))
		return nil, status.Error(codes.Internal, "系统繁忙")
	}

	// 类型转换处理
	result, ok := val.(int64)
	if !ok {
		l.Logger.Errorw("脚本返回类型异常",
			logx.Field("result", val),
			logx.Field("type", fmt.Sprintf("%T", val)))
		return nil, status.Error(codes.Internal, "系统异常")
	}

	// 处理结果
	switch result {
	case 0: // 扣减成功
		return &inventory.InventoryResp{}, nil
	case 1: // 已处理过
		l.Logger.Infow("订单已处理",
			logx.Field("pre_order_id", in.PreOrderId))
		resp.StatusCode = code.OrderhasBeenPaid
		resp.StatusMsg = code.OrderhasBeenPaidMsg
		return resp, nil

	case 2: // 库存不足
		resp.StatusCode = code.InventoryNotEnough
		resp.StatusMsg = code.InventoryNotEnoughMsg
		l.Logger.Infow("库存不足",
			logx.Field("pre_order_id", in.PreOrderId))
		return resp, nil
	default:
		l.Logger.Errorw("未知返回码",
			logx.Field("result", result))
		return nil, status.Error(codes.Internal, "系统异常")
	}
}
