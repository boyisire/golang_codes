package util

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"core/configs"
	pb "core/jobrequest"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	_ "github.com/joho/godotenv/autoload"
)

var (
	_kpOnce sync.Once
	_kp     *kproducer
)

type kproducer struct {
	Cfg      *configs.KafkaMqConfig
	Producer sarama.SyncProducer
}

func KafkaProducer(jsonCfgFile string) *kproducer {
	_kpOnce.Do(func() {
		cfg := &configs.KafkaMqConfig{}
		configs.LoadJsonConfig(cfg, jsonCfgFile)
		mqConfig := sarama.NewConfig()
		mqConfig.Net.SASL.Enable = true
		mqConfig.Net.SASL.User = cfg.Ak
		mqConfig.Net.SASL.Password = cfg.Password
		mqConfig.Net.SASL.Handshake = true
		certBytes, err := ioutil.ReadFile(cfg.FullCertFile)
		clientCertPool := x509.NewCertPool()
		ok := clientCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			log.Error("kafka producer failed to parse root certificate")
			return
		}
		mqConfig.Net.TLS.Config = &tls.Config{
			RootCAs:            clientCertPool,
			InsecureSkipVerify: true,
		}
		mqConfig.Net.TLS.Enable = true
		mqConfig.Producer.Return.Successes = true
		if err = mqConfig.Validate(); err != nil {
			msg := fmt.Sprintf("Kafka producer config invalidate. config: %v. err: %v", *cfg, err)
			log.Error(msg)
			return
		}
		producer, err := sarama.NewSyncProducer(cfg.Servers, mqConfig)
		if err != nil {
			msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
			log.Error(msg)
			return
		}
		_kp = &kproducer{
			Cfg:      cfg,
			Producer: producer,
		}
	})
	return _kp
}

func (me *kproducer) Publish(topic, key string, data *pb.JobRequest) error {
	if me.Producer == nil {
		log.WithFields(log.Fields{
			"topic": topic,
			"key":   key,
		}).Error("send kafka message error")
		return errors.New("kafka producer pointer nil")
	}
	content, _ := proto.Marshal(data)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(content),
	}
	_, _, err := me.Producer.SendMessage(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"topic": topic,
			"key":   key,
		}).Error("send kafka message error")
		return err
	}
	return nil
}
