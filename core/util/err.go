package util

import (
	log "github.com/Sirupsen/logrus"
)

// CommonError 微信返回的通用错误json
type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func FailOnError(err error, msg string) bool {
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"msg": msg,
		}).Error("Error receive data to connection!")
		return true
	}
	return false
}
