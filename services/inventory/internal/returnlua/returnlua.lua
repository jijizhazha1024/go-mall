
        -- 幂等性检查
        if redis.call("EXISTS", KEYS[1]) ~= 1 then
            return 1
        end  
        
        -- 归还库存
        for i=2, #KEYS do
            redis.call("INCRBY", KEYS[i], tonumber(ARGV[i]))
        end
        
       --删除锁
      	redis.call("DEL", KEYS[1])
	            
        return 0
    