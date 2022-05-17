package redis_store

const (
	CounterAlg = iota
	TokenBucketAlg
	LeakyBucketAlg
)

const CounterScript = `
	local key_prefix= KEYS[1]
     -- 微秒单位
    local unit = tonumber(ARGV[1])
    local max_count= tonumber(ARGV[2])
    local timestamp = redis.call("TIME")
    local second_time=tonumber(timestamp[1])
    local microsecond_time = tonumber(second_time*1000000+tonumber(timestamp[2]))
    local key=key_prefix..":"..math.floor(microsecond_time/unit)
    local n = redis.call("GET",key)
    local is_expire=0
    local expire=1
    if n == false then
       n=0
    else 
      is_expire=1
      n =tonumber(n)
    end 
    
    if n > max_count then 
     return 0
    end 
    
    local increment = max_count-n
    
    if increment <=0 then
      return 0
    end 

    redis.replicate_commands();
    redis.call("INCRBY", key, 1)
	if  is_expire==0  then
       
      if unit>1000000 then  
        expire =math.ceil(unit/1000000) 
      end

      redis.call("EXPIRE", key, expire)
    end 
	
    return increment
`

const TokenBucketScript = `
local key= KEYS[1]
     -- 微秒单位
local unit = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local capacity= tonumber(ARGV[3])
local timestamp = redis.call("TIME")
local second_time=tonumber(timestamp[1])
local microsecond_time = tonumber(second_time*1000000+tonumber(timestamp[2]))

local recentPutTokenTime = redis.call("HGET", key, "recentPutTokenTime")
if recentPutTokenTime == false then
  recentPutTokenTime=0;
else
 recentPutTokenTime=tonumber(recentPutTokenTime)
end

local distance_time= microsecond_time - recentPutTokenTime
local increment =math.floor(distance_time/unit*rate)

local num = redis.call("HGET", key, "num")

if num == false then
   num=0
else 
  
 num=tonumber(num)
end

num = math.min(num + increment, capacity)
num=math.ceil(num)
 redis.replicate_commands();
 if num >0 then
   if increment>0 then
      redis.call("HSET", key, "recentPutTokenTime", microsecond_time)
   end
   redis.call("HSET", key, "num", num-1)
 end
return num
`
