package user_coupons

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

type Status struct {
	ID     uint64 `db:"id"`
	Status int64  `db:"status"`
}
