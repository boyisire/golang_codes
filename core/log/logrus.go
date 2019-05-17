package log

import (
	"os"
	"reflect"

	log "github.com/Sirupsen/logrus"
)

type Fields = log.Fields
type Data = log.Fields

func init() {
	//处理日志格式
	//默认：log.SetFormatter(&log.TextFormatter{})
	//样式设置参考：https://codeday.me/bug/20190108/484553.html
	if os.Getenv("LOG_FMT") == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	//logrus.SetOutput(os.Stdout)

	//处理日志级别
	setLogLevel()
}

//设置日志级别
func setLogLevel() {
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

///////////////////日志输出规范化处理/////////////////

//格式化输出数据
func fmtFields(args ...interface{}) Fields {
	fields := make(Fields)
	if len(args) > 0 && reflect.TypeOf(args[0]) == reflect.TypeOf(fields) {
		fields = reflect.ValueOf(args[0]).Interface().(Fields)
	}
	return fields
}

func Debug(msg string, args ...interface{}) {
	log.WithFields(fmtFields(args...)).Debug(msg)
}

func Info(msg string, args ...interface{}) {
	log.WithFields(fmtFields(args...)).Info(msg)
}

func Warn(msg string, args ...interface{}) {
	log.WithFields(fmtFields(args...)).Warn(msg)
}

func Warning(msg string, args ...interface{}) {
	log.WithFields(fmtFields(args...)).Warning(msg)
}

func Error(msg string, args ...interface{}) {
	log.WithFields(fmtFields(args...)).Error(msg)
}

func Fatal(msg string, args ...interface{}) {
	log.WithFields(fmtFields(args...)).Fatal(msg)
}

func Panic(msg string, args ...interface{}) {
	log.WithFields(fmtFields(args...)).Panic(msg)
}

func NewLog() *log.Logger {
	return log.New()
}
