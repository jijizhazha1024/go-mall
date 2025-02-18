--  优惠券锁定脚本.lua

local userCouponKey = KEYS[1]
local preOrderKey = KEYS[2]
local couponId = ARGV[1]
local expireSeconds = ARGV[2]
local timestamp = ARGV[3]

-- 检查优惠券是否可用（原子化校验）
if redis.call("EXISTS", userCouponKey) == 1 then
    return 1
end

-- 存储锁定信息
redis.call("HSET", preOrderKey, couponId, timestamp)
redis.call("EXPIRE", preOrderKey, expireSeconds)
redis.call("SET", userCouponKey, "1", "EX", expireSeconds)
return 0
