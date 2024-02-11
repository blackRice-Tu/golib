[toc]

# GOLIB

封装若干golang常用模块，为后台开发和运维提效。


# 工具库介绍

## utils
通用工具库。

## xcache
本地缓存库

### xfreecache
基于开源库[freecache](https://github.com/coocood/freecache)的封装，具有以下特性：
* `freecache`原生特性。
* 支持缓存分组。

## xgin
HTTP框架库，基于开源库[gin](https://github.com/gin-gonic/gin)的封装，具有以下特性：
* 提供若干预设的错误码、返回结构。
* 提供`trace`中间件，为请求注入`trace_id`。
* 提供`cors`中间件。
* 提供`collect_request`中间件，结合xlogger，打印所有的http请求，允许从上下文自定义采集信息。
* 提供`recovery`中间件，从panic中恢复，并将异常堆栈写入logger。

## xgrpc
GRPC库，基于开源库[grpc-go](https://github.com/grpc/grpc-go)和[go-grpc-middleware](github.com/grpc-ecosystem/go-grpc-middleware)封装，具有以下特性：
* grpc client和server。
* 支持多服务间`trace_id`透传。
* 提供`collect_request`中间件，结合xlogger，打印所有的grpc请求，包括trace_id、阶段、grpc方法、grpc服务、请求内容、返回内容、状态码。
* 提供`recovery`中间件，从panic中恢复，并将异常堆栈写入logger。

## xkafka
kafka库，基于开源库[sarama](https://github.com/IBM/sarama)的封装，具有以下特性：
* kafka同步/异步生产者客户端。
* kafka消费者组客户端。
* debug模式下将kafka调用消息输出至终端。
* 支持`LoadOrNew`模式，将初始化后的producer实例写入全局缓存，多个producer实例以id区分。

## xlogger
日志库，基于开源库[zap](https://github.com/uber-go/zap)的封装，具有以下特性：
* 集成了[lumberjack](https://github.com/natefinch/lumberjack)，支持日志滚动切割。
* debug模式下将日志信息输出至终端。
* 支持`LoadOrNew`模式，将初始化后的logger实例写入全局缓存，多个logger实例以id区分。
* 自带默认logger实例，可直接import使用，支持`trace_id`追踪。
* 与xkafka配合使用，可将该logger实例下输出的所有日志收集到kafka指定的topic。

## xmysql
MySQL连接库，基于开源库[gorm](https://github.com/go-gorm/gorm)的封装，具有以下特性：
* MySQL连接客户端。
* 支持`LoadOrNew`模式，将初始化后的db实例写入全局缓存，多个db实例以id区分。
* 结合xlogger，打印所有SQL操作，包括trace_id、db实例id、操作类型、SQL语句、影响行数、耗时、异常信息。
* 提供`query_builder`和`order_builder`，支持复杂SQL查询语法。

## xredis
Redis连接库，基于开源库[go-redis](https://github.com/redis/go-redis)的封装，具有以下特性：
* Redis连接客户端。
* 支持`LoadOrNew`模式，将初始化后的db实例写入全局缓存，多个db实例以id区分。
* 可指定全局prefix，默认所有key都会带上prefix，也支持原始key操作。
* 提供若干工具，如`分布式限频器`，`分布式锁`。

## xrequests
HTTP连接库，基于开源库[grequests](https://github.com/levigross/grequests)的封装，具有以下特性：
* GET请求支持多种格式的query参数。
* POST请求支持多种格式的body参数。
* 所有请求都返回Trace，记录了URL、Method、Headers、Query、Body、StatusCode、ResponseBody信息。
* 结合xlogger，在debug模式下，或请求异常时，自动将Trace输出到logger。

## xtask
分布式任务库，基于开源库[async](https://github.com/hibiken/asynq)的封装，具有以下特性：
* 结合xredis库，将redis作为broker。
* 支持分布式定时任务、异步任务。

# 使用说明

## 后台服务
后台服务，推荐按照以下顺序使用：
1. 生成本地配置文件，填写默认日志、HTTP日志、MySQL日志、GRPC日志信息，使用开源`viper`库读取配置。
2. 使用`golib.Option`生成golib配置，并通过`golib.SetConfig(opt)`注册配置。
3. 使用`xkafka`创建kafka生产者。
4. 使用`xlogger`初始化默认logger，并将日志hook到已创建的kafka生产者中。
5. 使用`xmysql`初始化MySQL连接。
6. 使用`xredis`初始化Redis连接。
7. 使用`xtask`初始化分布式任务服务。
8. 使用`xfreecache`初始化本地缓存，并从MySQL或Redis中将数据加载到缓存。
9. 使用`gin`初始化HTTP服务，并注册`xgin`中的中间件。
10. 使用`grpc-go`初始化GRPC服务，并注册`xgrpc`中的中间件。
11. 启动服务。
12. 从kafka中收集日志，入库，监控告警。
