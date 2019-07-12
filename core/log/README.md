## 功能描述
该库的主要功能是，整合日志处理,方便项目统一管理。

**目标：** 实现一个高可用的日志管理系统

## 功能点
目前主要实现的功能点：
+ 格式化日志输出样式
+ 格式化日志级别
+ 格式化日志记录方法
+ 日志统一发送ES
>es_index: logs

+ 错误日志报警
>通过worker-log-biz 异步处理，并发送钉钉通知

## 使用样式
```
引用：
"core/log"

使用：
1. 不带参打印日志信息
log.Info("doWorker Begin...")

2. 带参调试日志
log.Debug("go func doSubJob", log.Data{
    "i":      i,
    "j":      j,
})

3. 带参错误日志
log.Error("error msg", log.Data{
    "err":   err,
})

4. 原始用法
log.New().WithFields(log.Fields{
    "a": 1,
    "b": 2,
}).Info("test")

.env设置：
#调试开关(true/false)
DEBUG=false
#日志级别
LOG_LEVEL=debug
#日志样式
LOG_FMT=json
......
```


## 引用的核心库
+ [logrus](http://github.com/Sirupsen/logrus)

## 常用配置

+ DEBUG `调试开关` *(默认：false)*

+ LOG_FMT `日志格式` *(默认：text)*
  - text
  - json

+ LOG_LEVEL  `显示日志级别` *(默认：info)*

+ LOG_TO_ES  `是否写入ES` *(默认：true)*

+ LOG_ALARM_NOTIFY  `是否发送报警` *(默认：true)*
