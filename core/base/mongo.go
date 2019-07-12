package base

import (
	"errors"
	"fmt"
	"os"
	"time"

	"core/log"
	_ "github.com/joho/godotenv/autoload"
	mgo "gopkg.in/mgo.v2"
)

var _mongo *Mongo

type Mongo struct {
	Session  *mgo.Session
	DataBase string
}

/**
 * 获取mongo session
 */
func NewMongo() *Mongo {
	if _mongo == nil {
		_mongo = &Mongo{}
	}

	return _mongo
}

func (me *Mongo) Open(mongo_type string) error {
	switch mongo_type {
	case "default":
		me.Session = getSession(os.Getenv("MONGO_URL"))
		me.DataBase = os.Getenv("MONGO_DATABASE")
	default:
		return errors.New("MongoDB Type Error!")
	}
	return nil
}

func getSession(url string) *mgo.Session {
	session, err := mgo.Dial(url)
	if err != nil {
		time.Sleep(time.Duration(retries) * time.Second)
		log.Error(fmt.Sprintf("DB-mongo Connection Fail:: err=[%s], retries=[%d]", err.Error(), retries))
		retries += 2
		getSession(url)
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}
