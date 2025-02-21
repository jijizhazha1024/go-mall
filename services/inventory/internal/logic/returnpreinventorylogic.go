package logic

import (
	"context"
	"fmt"

	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReturnPreInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReturnPreInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReturnPreInventoryLogic {
	return &ReturnPreInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReturnPreInventory 退还预扣减的库存（）
func (l *ReturnPreInventoryLogic) ReturnPreInventory(in *inventory.InventoryReq) (*inventory.InventoryResp, error) {

	// 构建幂等锁Key（用户ID+预订单ID）
	lockKey := fmt.Sprintf("inventory:return:lock:%d:%s", in.UserId, in.PreOrderId)

	// 准备Lua脚本参数
	keys := []string{lockKey}
	args := []interface{}{in.PreOrderId}
	productKeys := make([]string, 0, len(in.Items))

	// 构造库存Key列表
	for _, item := range in.Items {
		if item.Quantity <= 0 {
			return nil, status.Error(codes.InvalidArgument, "商品数量不合法")
		}
		productKey := fmt.Sprintf("inventory:product:%d", item.ProductId)
		productKeys = append(productKeys, productKey)
		args = append(args, item.Quantity)
	}
	keys = append(keys, productKeys...)

	// 组装Lua脚本
	luaScript := `
        -- 幂等性检查
        if redis.call("EXISTS", KEYS[1]) == 1 then
            return 1
        end
        
        
        -- 归还库存
        for i=2, #KEYS do
            redis.call("INCRBY", KEYS[i], tonumber(ARGV[i]))
        end
        
        -- 设置处理标记（30分钟过期）
        redis.call("SET", KEYS[1], ARGV[1], "EX", 1800)
        return 0
    `

	// 执行Lua脚本（使用go-zero的Eval方法）
	val, err := l.svcCtx.Rdb.Eval(luaScript, keys, args)
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
	case 0: // 归还成功
		return &inventory.InventoryResp{}, nil
	case 1: // 已处理过
		return &inventory.InventoryResp{}, status.Error(codes.AlreadyExists, "订单已处理")

	default:
		l.Logger.Errorw("未知返回码",
			logx.Field("result", result))
		return nil, status.Error(codes.Internal, "系统异常")
	}

}
