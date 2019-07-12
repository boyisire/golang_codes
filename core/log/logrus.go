package log

import (
	"encoding/json"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"core/util"
	"github.com/Sirupsen/logrus"
)

type (
	Fields = logrus.Fields
	Data   = logrus.Fields
	Entry  = logrus.Entry
	Level  = logrus.Level
)

var (
	Logger         *logrus.Logger
	Log            *logrus.Entry
	Hostname, _    = os.Hostname()
	LogToEs        = true
	LogAlarmNotify = true
	ThreadId       = 0
	ClientIp       = "0.0.0.0"
	AppName        = ""
)

func init() {
	logger := NewLogger()
	// 输出类型:: 设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logger.Out = os.Stdout
	//设置日志格式
	if os.Getenv("LOG_FMT") == "json" {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	//处理日志级别
	setLogLevel(logger)

	//初始化日志输出格式
	now := time.Now()
	Log = logrus.NewEntry(logger)
	Log = Log.WithFields(Fields{
		"hostname":    Hostname,    //主机名
		"pid":         os.Getpid(), //进程号
		"timestamp":   now.Unix(),  //时间戳
		"microsecond": 0,           //微秒
		"thread_id":   ThreadId,    //事务id
		"app_name":    AppName,     //服务名称址
		"file":        "",          //执行文件
		"func":        "",          //执行方法
		"line":        0,           //错误行号
		"msg":         "",          //打印消息
		"extends":     "",          //打印数据
		"created_at":  now.Format("2006-01-02 15:04:05"),
	})
	// 写入ES控制
	if os.Getenv("LOG_TO_ES") == "false" {
		LogToEs = false
	}
}

// logrus New
func New() *logrus.Logger {
	return logrus.New()
}

// 新建Logger
func NewLogger() *logrus.Logger {
	if Logger != nil {
		return Logger
	}
	Logger = New()
	return Logger
}

// 设置日志级别
func setLogLevel(logger *logrus.Logger) {
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logger.Level = logrus.DebugLevel
	case "info":
		logger.Level = logrus.InfoLevel
	case "warning":
		logger.Level = logrus.WarnLevel
	case "error":
		logger.Level = logrus.ErrorLevel
	case "fatal":
		logger.Level = logrus.FatalLevel
	case "panic":
		logger.Level = logrus.PanicLevel
	default:
		logger.Level = logrus.InfoLevel

	}
}

// 推送到Amq
func PushLogToAmq(fields Fields) {
	logStr, _ := json.Marshal(fields)
	util.AmqpPublish2(os.Getenv("MQ_URL"), "log", &util.AmqpJobData{
		Command: "app-log-collect",
		Data: map[string]string{
			"logs": string(logStr),
		},
	})
}

// 生产事务流水号
func BuildThreadId() int {
	ThreadId, _ = strconv.Atoi(util.RandCode(10))
	return ThreadId
}

///////////////////日志格式化处理/////////////////

// 日志格式处理
func logPut(level string, msg string, args ...interface{}) {
	fields := make(Fields)
	if len(args) > 0 && reflect.TypeOf(args[0]) == reflect.TypeOf(fields) {
		fields = reflect.ValueOf(args[0]).Interface().(Fields)
		if thread_id, ok := fields["thread_id"]; ok {
			Log.Data["thread_id"] = thread_id
			delete(fields, "thread_id")
		}
		// 客户端IP
		if client_ip, ok := fields["client_ip"]; ok {
			Log.Data["client_ip"] = client_ip
			delete(fields, "client_ip")
		}
		// 是否写入ES控制
		if log_to_es, ok := fields["log_to_es"]; ok {
			if log_to_es == false {
				LogToEs = false
			}
			delete(fields, "log_to_es")
		}
		// 是否发送日志报警:: 调试期间，方了减少报警，可直接在输出数据中，写：(log_alarm_notify:false) 关闭报警
		if log_alarm_notify, ok := fields["log_alarm_notify"]; ok {
			if log_alarm_notify == false {
				LogAlarmNotify = false
			}
			delete(fields, "log_alarm_notify")
		}
	}
	Log.Data["created_at"] = time.Now().Format("2006-01-02 15:04:05")

	extends, _ := json.Marshal(fields)
	Log.Data["extends"] = string(extends)
	if ThreadId != 0 && Log.Data["thread_id"] == 0 {
		Log.Data["thread_id"] = ThreadId
	}
	if ClientIp != "0.0.0.0" && Log.Data["client_ip"] == "" {
		Log.Data["client_ip"] = ClientIp
	}
	//当前时间戳
	timestamp := time.Now().Unix()
	Log.Data["timestamp"] = timestamp
	//当前微秒
	microsecond := strconv.FormatInt(time.Now().UnixNano()/1000, 10)[10:]
	Log.Data["microsecond"], _ = strconv.Atoi(microsecond)
	//调用文件信息，向上倒2级
	if pc, file, line, ok := runtime.Caller(2); ok {
		Log.Data["file"] = path.Base(file)
		Log.Data["line"] = line
		Log.Data["func"] = runtime.FuncForPC(pc).Name()
		if AppName != "" && Log.Data["app_name"] == "" {
			Log.Data["app_name"] = AppName
		} else {
			Log.Data["app_name"] = path.Base(path.Dir(file))
		}
	}

	//通过队列,发送到ES
	if LogToEs && (os.Getenv("LOG_LEVEL") == "debug" || level != "debug") {
		//格式化数据，方便统一ES数据格式
		logData := make(Fields)
		logData["app_name"] = Log.Data["app_name"]
		logData["hostname"] = Log.Data["hostname"]
		logData["pid"] = Log.Data["pid"]
		logData["timestamp"] = Log.Data["timestamp"]
		logData["microsecond"] = Log.Data["microsecond"]
		logData["thread_id"] = Log.Data["thread_id"]
		logData["file"] = Log.Data["file"]
		logData["func"] = Log.Data["func"]
		logData["line"] = Log.Data["line"]
		logData["level"] = level
		logData["msg"] = msg
		logData["extends"] = Log.Data["extends"]
		logData["created_at"] = Log.Data["created_at"]
		logData["log_alarm_notify"] = LogAlarmNotify
		if _, ok := Log.Data["client_ip"]; ok {
			logData["client_ip"] = Log.Data["client_ip"]
		}
		PushLogToAmq(logData)
	}

	// 日志输出
	switch level {
	case "debug":
		Log.Debug(msg)
	case "info":
		Log.Info(msg)
	case "warn":
		Log.Warn(msg)
	case "warning":
		Log.Warning(msg)
	case "error":
		Log.Error(msg)
	case "fatal":
		Log.Fatal(msg)
	case "panic":
		Log.Panic(msg)
	default:
		Log.Info(msg)
	}
	return
}

func Debug(msg string, args ...interface{}) {
	logPut("debug", msg, args...)
}

func Info(msg string, args ...interface{}) {
	logPut("info", msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logPut("warn", msg, args...)
}

func Warning(msg string, args ...interface{}) {
	logPut("warn", msg, args...)
}

func Error(msg string, args ...interface{}) {
	logPut("error", msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logPut("fatal", msg, args...)
}

func Panic(msg string, args ...interface{}) {
	logPut("panic", msg, args...)
}
