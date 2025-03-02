
        -- 幂等性检查
        if redis.call("EXISTS", KEYS[1]) ~= 1 then
            return 1
        end  
        
        -- 归还库存
        for i=2, #KEYS do
            local exists = redis.call("EXISTS", KEYS[i])
            if exists == 1 then
            redis.call("INCRBY", KEYS[i], tonumber(ARGV[i]))
        end
    end
        
       --删除锁
      	redis.call("DEL", KEYS[1])
	            
        return 0
    