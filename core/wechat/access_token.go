package wechat

import (
	"fmt"

	"core/util"

	"github.com/garyburd/redigo/redis"
)

// ResAccessToken struct
type ResAccessToken struct {
	util.CommonError
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// GetAccessToken 获取access_token
func (me *Wechat) GetAccessToken(appid string) (accessToken string, err error) {
	redisConn := me.RedisPool.Get()
	redisConn.Do("SELECT", 0)
	me.accessTokenLock.Lock()
	defer func() {
		me.accessTokenLock.Unlock()
		redisConn.Close()
	}()

	accessTokenCacheKey := fmt.Sprintf("17jzh.wechat.access_token.%s", appid)
	accessToken, err = redis.String(redisConn.Do("GET", accessTokenCacheKey))
	return
}
