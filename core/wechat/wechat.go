package wechat

import (
	"strconv"
	"sync"
	"time"

	"17jzh.com/core/base"
	"17jzh.com/core/util"

	"github.com/garyburd/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
)

var _wechat *Wechat

type Wechat struct {
	RedisPool     *redis.Pool
	WechatAccount *WechatAccount
	//accessTokenLock 读写锁 同一个AppID一个
	accessTokenLock *sync.RWMutex
}

type WechatAccount struct {
	ID     int
	Appid  string
	Secret string
	Token  string
	Mcid   string
	Mckey  string
}

func NewWechat(ctx *base.Context) *Wechat {
	if _wechat == nil {
		_wechat = &Wechat{RedisPool: ctx.Redis.Default}
		_wechat.setAccessTokenLock(new(sync.RWMutex))
	}
	return _wechat
}

//SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (me *Wechat) setAccessTokenLock(lock *sync.RWMutex) {
	me.accessTokenLock = lock
}
func (me *Wechat) SetWechatAccount(wechatAccount *WechatAccount) {
	me.WechatAccount = wechatAccount
}

func (me *Wechat) addPushErrRedisLog(memberKey string, err error, errCode int64, userID int) {
	if err == nil {
		return
	}
	conn := me.RedisPool.Get()
	conn.Do("SELECT", 1) //选择数据库1
	defer conn.Close()

	util.FailOnError(err, strconv.Itoa(userID))
	conn.Do("ZADD", memberKey, errCode, userID)
}

func (me *Wechat) addPushNewsSuccessLog(memberKey string, userID int) {
	conn := me.RedisPool.Get()
	conn.Do("SELECT", 1) //选择数据库1
	defer conn.Close()

	conn.Do("ZADD", memberKey, time.Now().Unix(), userID)
}
