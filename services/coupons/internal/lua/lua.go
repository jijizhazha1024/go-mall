package lua

import _ "embed"

//go:embed lock_coupon.lua
var LockCouponScript string

//go:embed unlock_coupon.lua
var UnlockCouponScript string
