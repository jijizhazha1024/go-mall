package response

import (
	"jijizhazha1024/go-mall/common/consts/code"
)

type RefreshResponse struct {
	Response
	Data interface{} `json:"data"`
}

func NewRefreshResponse(data interface{}) RefreshResponse {
	return RefreshResponse{
		Response: Response{
			StatusCode: code.TokenRenewed,
			StatusMsg:  code.TokenRenewedMsg,
		},
		Data: data,
	}
}
