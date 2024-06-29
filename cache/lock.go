package cache

import "github.com/bsm/redislock"

var Locker *redislock.Client

const (
	UserEditProfileLock = "user_profile_lock"
)

func UserEditProfileLockKey(userName string) string {
	return UserEditProfileLock + userName
}
