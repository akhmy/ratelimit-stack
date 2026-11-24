local key = KEYS[1]
if not key then
    return redis.error_reply('ERR invalid key')
end

local window = tonumber(ARGV[1])
if not window or window <= 0 or window % 1 ~= 0 then
    return redis.error_reply('ERR invalid window')
end

local limit = tonumber(ARGV[2])
if not limit or limit <= 0 or limit % 1 ~= 0 then
    return redis.error_reply('ERR invalid limit')
end

local time = tonumber(ARGV[3])
if not time or time < 0 or time % 1 ~= 0 then
    return redis.error_reply('ERR invalid time')
end

local requestID = ARGV[4]
if not requestID or requestID == "" then
    return redis.error_reply('ERR invalid requestID')
end

redis.call('ZREMRANGEBYSCORE', key, '-inf', time - window)

if redis.call('ZCARD', key) >= limit then
    return 0
end

redis.call('ZADD', key, time, requestID)
redis.call('PEXPIRE', key, window)
return 1
