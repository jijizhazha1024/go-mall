package logic

import (
	"jijizhazha1024/go-mall/apis/carts/internal/types"
	"jijizhazha1024/go-mall/services/carts/carts"
)

func ConvertCartInfoResponse(res []*carts.CartInfoResponse) []*types.CartInfoResponse {
	var result []*types.CartInfoResponse
	for _, item := range res {
		result = append(result, &types.CartInfoResponse{
			Id:        item.Id,
			UserId:    item.UserId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Checked:   item.Checked,
		})
	}
	return result
}
