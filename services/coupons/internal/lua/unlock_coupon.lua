-- 释放优惠券原子操作
local userCouponKey = KEYS[1]
local preOrderKey = KEYS[2]
local couponId = ARGV[1]

-- 原子化释放操作
if redis.call("EXISTS", userCouponKey) == 0 then
    return 1 -- 表示优惠券未锁定
end

redis.call("DEL", userCouponKey)
redis.call("HDEL", preOrderKey, couponId)
return 0 -- 成功释放