package base

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"core/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type gormLogger struct{}

var _db *DB

// 打印SQL调试语句
func (*gormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		log.Debug("SQL", log.Fields{
			"query": v[3],
			"bind":  v[4],
			"time":  v[2],
		})
	}
}

// DB database manager class
type DB struct {
	*gorm.DB
}

// NewDB returns a new Model without opening database connection
func NewDB() *DB {
	if _db == nil {
		_db = &DB{}
	}
	return _db
}

func (me *DB) Close() {
	me.Close()
}

// Open opens database connection with the settings found in cfg
func (me *DB) Open(db_type string) error {
	switch db_type {
	case "default":
		me.Conn("DB_URL")
		return nil
	default:
		return errors.New("DB Type Error")
	}
}

// DB连接
func (me *DB) Conn(dbConnName string) *DB {
	dbs := os.Getenv(dbConnName)
	if dbs == "" {
		dbs = os.Getenv("DB_URL")
	}
	//创建连接句柄
	conn := connDb(dbs)

	//连接池处理::
	// SetMaxIdleConns 设置空闲连接池中的最大连接数。
	_dbMaxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MaxIdleConns"))
	if _dbMaxIdleConns == 0 {
		_dbMaxIdleConns = 10
	}
	conn.DB().SetMaxIdleConns(_dbMaxIdleConns)

	// SetMaxOpenConns 设置数据库连接最大打开数。
	_dbMaxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MaxOpenConns"))
	if _dbMaxOpenConns == 0 {
		_dbMaxOpenConns = 100
	}
	conn.DB().SetMaxOpenConns(_dbMaxOpenConns)

	// SetConnMaxLifetime 设置可重用连接的最长时间
	_dbMaxLifetime, _ := strconv.Atoi(os.Getenv("DB_Maxlifetiem"))
	if _dbMaxLifetime == 0 {
		_dbMaxLifetime = 14400
	}
	conn.DB().SetConnMaxLifetime(time.Duration(_dbMaxLifetime) * time.Second)

	// 调试模式
	if debug == "true" {
		conn.LogMode(true)
		conn.SetLogger(&gormLogger{})
	}

	me.DB = conn
	return me
}

// DB连接句柄
func connDb(dbs string) *gorm.DB {
	conn, err := gorm.Open("mysql", dbs)
	if err != nil {
		time.Sleep(time.Duration(retries) * time.Second)
		log.Error(fmt.Sprintf("DB-mysql Connection Fail:: err=[%s], retries=[%d]", err.Error(), retries))
		retries += 2
		connDb(dbs)
	}
	return conn
}
