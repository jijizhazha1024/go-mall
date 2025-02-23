package user_coupons

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/services/coupons/coupons"
)

var _ UserCouponsModel = (*customUserCouponsModel)(nil)

type (
	// UserCouponsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserCouponsModel.
	UserCouponsModel interface {
		userCouponsModel
		WithSession(session sqlx.Session) UserCouponsModel
		QueryUserCoupons(ctx context.Context, userId, page, pageSize int32) ([]*UserCoupons, error)
		CheckUserCouponExistWithLock(ctx context.Context, session sqlx.Session, userId uint64, couponId string) (bool, error)
		GetUserCouponByUserIdCouponIdWithLock(ctx context.Context, session sqlx.Session, userId uint64, couponId string) (*UserCoupons, error)
		GetStatusByUserIdCouponId(ctx context.Context, userid int32, couponId string) (*Status, error)
		UpdateStatusOrderById(ctx context.Context, orderId string, id int, status coupons.CouponUsageStatus) error
	}

	customUserCouponsModel struct {
		*defaultUserCouponsModel
	}
)

func (m *customUserCouponsModel) UpdateStatusOrderById(ctx context.Context, orderId string, id int, used coupons.CouponUsageStatus) error {
	query := fmt.Sprintf("update %s set `status` = ?, `order_id` = ?,used_at = now() where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, used, orderId, id)
	return err
}

func (m *customUserCouponsModel) GetStatusByUserIdCouponId(ctx context.Context, userid int32, couponId string) (*Status, error) {
	var status Status
	query := fmt.Sprintf("select `id`,`status` from %s where `user_id` = ? and `coupon_id` = ? limit 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &status, query, userid, couponId)
	switch {
	case err == nil:
		return &status, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserCouponsModel) GetUserCouponByUserIdCouponIdWithLock(ctx context.Context, session sqlx.Session, userId uint64, couponId string) (*UserCoupons, error) {
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

func (m *customUserCouponsModel) CheckUserCouponExistWithLock(ctx context.Context, session sqlx.Session, userId uint64, couponId string) (bool, error) {
	var cnt int64
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `coupon_id` = ? LIMIT 1 FOR SHARE ", m.table)
	err := session.QueryRowCtx(ctx, &cnt, query, userId, couponId)
	switch {
	case err == nil:
		return cnt > 0, nil
	default:
		return false, err
	}
}

func (m *customUserCouponsModel) QueryUserCoupons(ctx context.Context, userId, page, pageSize int32) ([]*UserCoupons, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by `created_at` desc limit ?,?", userCouponsRows, m.table)
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
