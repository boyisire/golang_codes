package base

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"core/configs"
	pb "core/jobrequest"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/golang/protobuf/proto"
)

type Kafka struct {
	Cfg      *configs.KafkaMqConfig
	Consumer *cluster.Consumer
	Sig      chan os.Signal
}

func NewKafka(jsonCfgFile string) *Kafka {
	//load and set config
	cfg := &configs.KafkaMqConfig{}
	configs.LoadJsonConfig(cfg, jsonCfgFile)
	clusterCfg := cluster.NewConfig()
	clusterCfg.Net.SASL.Enable = true
	clusterCfg.Net.SASL.User = cfg.Ak
	clusterCfg.Net.SASL.Password = cfg.Password
	clusterCfg.Net.SASL.Handshake = true
	certBytes, err := ioutil.ReadFile(cfg.FullCertFile)
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("kafka consumer failed to parse root certificate")
	}
	clusterCfg.Net.TLS.Config = &tls.Config{
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}
	clusterCfg.Net.TLS.Enable = true
	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	clusterCfg.Group.Return.Notifications = true
	clusterCfg.Version = sarama.V0_10_0_0
	if err = clusterCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka consumer config invalidate. config: %v. err: %v", *clusterCfg, err)
		panic(msg)
	}
	//new consumer
	consumer, err := cluster.NewConsumer(cfg.Servers, cfg.ConsumerId, cfg.Topics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		panic(msg)
	}
	//new sig
	sig := make(chan os.Signal, 1)

	return &Kafka{
		Cfg:      cfg,
		Consumer: consumer,
		Sig:      sig,
	}
}

func (me *Kafka) Consume(job IJob, async bool) {
	for {
		select {
		case msg, more := <-me.Consumer.Messages():
			if more {
				log.Info("key : ", string(msg.Key))
				preq := &pb.JobRequest{}
				err := proto.Unmarshal(msg.Value, preq)
				if err != nil {
					log.Error(err)
				}
				if async == true {
					go func() {
						job.Handle(preq)
						me.Consumer.MarkOffset(msg, "") // mark message as processed
					}()
				} else {
					job.Handle(preq)
					me.Consumer.MarkOffset(msg, "") // mark message as processed
				}
			}
		case err, more := <-me.Consumer.Errors():
			if more {
				log.Error(err.Error())
			}
		case _, more := <-me.Consumer.Notifications():
			if more {
				log.Info("Consumer rebalance...")
			}
		case <-me.Sig:
			me.Consumer.Close()
			return
		}
	}
}

func (me *Kafka) Stop(s os.Signal) {
	me.Sig <- s
	log.Info("Stop consumer server!!!")
}
