package user_coupons

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserCouponsModel = (*customUserCouponsModel)(nil)

type (
	// UserCouponsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserCouponsModel.
	UserCouponsModel interface {
		userCouponsModel
		WithSession(session sqlx.Session) UserCouponsModel
		QueryUserCoupons(ctx context.Context, userId, page, pageSize int32) ([]*UserCoupons, error)
		CheckUserCouponExistWithSession(ctx context.Context, session sqlx.Session, userId uint64, couponId string) (bool, error)
		GetUserCouponByUserIdCouponIdLock(ctx context.Context, session sqlx.Session, userId uint64, couponId string) (*UserCoupons, error)
	}

	customUserCouponsModel struct {
		*defaultUserCouponsModel
	}
)

func (m *customUserCouponsModel) GetUserCouponByUserIdCouponIdLock(ctx context.Context, session sqlx.Session, userId uint64, couponId string) (*UserCoupons, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `coupon_id` = ? limit 1 for update", userCouponsRows, m.table)
	var resp UserCoupons
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, couponId)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserCouponsModel) CheckUserCouponExistWithSession(ctx context.Context, session sqlx.Session, userId uint64, couponId string) (bool, error) {
	var cnt int64
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `coupon_id` = ?", m.table)
	err := session.QueryRowCtx(ctx, &cnt, query, userId, couponId)
	switch {
	case err == nil:
		return cnt > 0, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (m *customUserCouponsModel) QueryUserCoupons(ctx context.Context, userId, page, pageSize int32) ([]*UserCoupons, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit ?,?", userCouponsRows, m.table)
	var resp []*UserCoupons
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, (page-1)*pageSize, pageSize)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return resp, nil
	default:
		return nil, err
	}
}

// NewUserCouponsModel returns a model for the database table.
func NewUserCouponsModel(conn sqlx.SqlConn) UserCouponsModel {
	return &customUserCouponsModel{
		defaultUserCouponsModel: newUserCouponsModel(conn),
	}
}

func (m *customUserCouponsModel) WithSession(session sqlx.Session) UserCouponsModel {
	return NewUserCouponsModel(sqlx.NewSqlConnFromSession(session))
}
