## 功能描述
该库的主要功能是，整合日志处理,方便项目统一管理。

**目标：** 实现一个高可用的日志管理系统

## 功能点
目前主要实现的功能点：
+ 格式化日志输出样式
+ 格式化日志级别
+ 格式化日志记录方法


## 使用样式
```
引用：
"app/core/log"

使用：
1.
log.Info("doWorker Begin...")

2.
log.Debug("go func doSubJob", log.Data{
    "i":      i,
    "j":      j,
})

3.
log.Error("error msg", log.Data{
    "err":   err,
})

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

