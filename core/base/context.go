package base

import (
	"os"
	"strings"
)

// Context wraps request and response. It provides methods for handling responses
type Context struct {
	//DB is the database stuff, with all models registered
	DB    *DB
	DBS   map[string]*DB
	Redis *Redis
	Mongo *Mongo
	Es    *Es
	Cfg   map[string]string
}

var (
	// 调试开关
	debug = os.Getenv("DEBUG")
	// 连接失败重试等待间隔,每次增2
	retries = 2
)

// NewContext creates new context for the given w and r
func NewContext(cfg map[string]string) (*Context, error) {
	ctx := &Context{}
	ctx.Cfg = cfg
	err := ctx.ctxInit()
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

// Init initializes the context
func (c *Context) ctxInit() error {
	if c.Cfg["db"] == "default" {
		db := NewDB()
		arr := strings.Split(c.Cfg["db"], ",")

		for _, v := range arr {
			if err := db.Open(strings.TrimSpace(v)); err != nil {
				return err
			}
		}
		c.DB = db
	}

	if c.Cfg["redis"] != "" {
		redis := NewRedis()
		arr := strings.Split(c.Cfg["redis"], ",")
		for _, v := range arr {
			if err := redis.Open(strings.TrimSpace(v)); err != nil {
				return err
			}
		}
		c.Redis = redis
	}

	if c.Cfg["mongo"] == "default" {
		mongo := NewMongo()
		if err := mongo.Open(strings.TrimSpace(c.Cfg["mongo"])); err != nil {
			return err
		}
		c.Mongo = mongo
	}

	if c.Cfg["es"] != "" {
		es := NewEs()
		if err := es.Open(strings.TrimSpace(c.Cfg["es"])); err != nil {
			return err
		}
		c.Es = es
	}
	return nil
}

// 多DB连接
func (c *Context) NewContextDbs(dbs map[int]string) error {
	db := NewDB()
	for _, dbSuff := range dbs {
		conn := db.Conn(dbSuff)
		c.DBS[dbSuff] = conn
	}

	return nil
}
