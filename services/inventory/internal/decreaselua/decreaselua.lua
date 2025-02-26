 -- 幂等性检查
 if redis.call("EXISTS", KEYS[1]) == 1 then
    return 1
end

-- 预检查库存
for i=2, #KEYS do
    local stock = tonumber(redis.call('GET',KEYS[i]) or 0)
    local deduct = tonumber(ARGV[i])  -- ARGV索引从1开始
    if stock < deduct then
    --删除锁
    redis.call("DEL", KEYS[1])
        return 2
    end
end

-- 扣减库存
for i=2, #KEYS do
redis.call('DECRBY', KEYS[i], tonumber(ARGV[i]))
end

-- 设置处理标记（30分钟过期）
redis.call("SET", KEYS[1], ARGV[1], "EX", 1800)
return 0