package coupon

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var ErrNotFound = sqlx.ErrNotFound

type CStatus struct {
	EndTime time.Time `db:"end_time"` // 有效期结束
	Status  int64     `db:"status"`   // 状态：0-禁用 1-启用
}
