package base

import (
	"errors"
	"fmt"
	"os"
	"time"

	"core/log"
	"core/util"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	elastic "gopkg.in/olivere/elastic.v5"
)

type Es struct {
	Client *elastic.Client
}

var (
	_es     *Es
	EsHosts map[int]string = map[int]string{
		0: os.Getenv("ES_HOST_MASTER"), //主
		1: os.Getenv("ES_HOST_SLAVE"),  //从
	}
)

func NewEs() *Es {
	if _es == nil {
		_es = &Es{}
	}
	return _es
}

func (me *Es) Open(esType string) error {
	switch esType {
	case "master":
		me.Conn(EsHosts[0])
		return nil
	case "slave":
		me.Conn(EsHosts[1])
		return nil
	default:
		return errors.New("ES Type Error")
	}
}

// DB连接
func (me *Es) Conn(esHost string) *Es {
	isReconnect := false
	client, err := connES(esHost, isReconnect)
	if client == nil && err != nil {
		log.Error(fmt.Sprintf("DB-Es-1 Connection Fail:: err=[%s], Host=[%s] ", err.Error(), esHost))
		if esHost == EsHosts[0] {
			esHost = EsHosts[1]
		} else {
			esHost = EsHosts[0]
		}
		client, err = connES(esHost, isReconnect)
		if client == nil && err != nil {
			log.Error(fmt.Sprintf("DB-Es-2 Connection Fail:: err=[%s], Host=[%s] ", err.Error(), esHost))
			//主从都不通的情况下,定期随机重试连接
			esHost = EsHosts[int(util.Random(1))]
			isReconnect = true
			client, err = connES(esHost, isReconnect)
		}
	}
	me.Client = client
	return me
}

// ES连接句柄
func connES(esHost string, isReconnect bool) (client *elastic.Client, err error) {
	client, err = elastic.NewClient(elastic.SetURL(esHost))
	if err != nil && isReconnect {
		time.Sleep(time.Duration(retries) * time.Second)
		log.Error(fmt.Sprintf("DB-Es-n Connection Fail:: err=[%s], retries=[%d], host=[%s]", err.Error(), retries, esHost))
		retries += 2
		connES(esHost, true)
	}
	return
}
