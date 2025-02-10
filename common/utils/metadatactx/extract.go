package metadatactx

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func ExtractFromMetadataCtx(ctx context.Context, key string) (string, bool) {
	// 获取客户端IP
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	values := md.Get(key)
	if len(values) == 0 {
		return "", false
	}
	return values[0], true
}
