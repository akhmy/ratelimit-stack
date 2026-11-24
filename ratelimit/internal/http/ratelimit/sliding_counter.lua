local cur_key = KEYS[1]
if not cur_key then
    return redis.error_reply('ERR invalid current key')
end

local prev_key = KEYS[2]
if not prev_key then
    return redis.error_reply('ERR invalid prev key')
end

local window_ms = tonumber(ARGV[1])
if not window_ms or window_ms <= 0 or window_ms % 1 ~= 0 then
    return redis.error_reply('ERR invalid window')
end

local limit = tonumber(ARGV[2])
if not limit or limit <= 0 or limit % 1 ~= 0 then
    return redis.error_reply('ERR invalid limit')
end

local unix_ms = tonumber(ARGV[3])
if not unix_ms or unix_ms < 0 or unix_ms % 1 ~= 0 then
    return redis.error_reply('ERR invalid time')
end

local cur_cnt = tonumber(redis.call('GET', cur_key) or '0')
local prev_cnt = tonumber(redis.call('GET', prev_key) or '0')

if not cur_cnt or cur_cnt < 0 or cur_cnt % 1 ~= 0
    or not prev_cnt or prev_cnt < 0 or prev_cnt % 1 ~= 0 then
    return redis.error_reply('ERR invalid window counter')
end

-- Calculating n of requests.
local elapsed_ms = unix_ms % window_ms
local prev_weight = (window_ms - elapsed_ms) / window_ms
local estimated_cnt = prev_weight * prev_cnt + cur_cnt + 1
if estimated_cnt > limit then
    return 0
end

redis.call('INCR', cur_key)
redis.call('PEXPIREAT', cur_key, unix_ms - elapsed_ms + 2 * window_ms)
return 1
