package redis

const (
	KeyBan      = "ban"
	KeyRequests = "requests"
	KeyBanLevel = "banlevel"

	KeyMinAmount = "min_amount"

	KeyEvents = "events"

	KeyMinAmountLock = "key_amount_lock"
)

const UnlockScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`
