// Package config 配置，所有配置会同步到这里
package config

import (
	"github.com/spf13/viper"
	"time"
)

// C 配置的单例
var C = defaultConfig()

// envConfig 配置结构体
type envConfig struct {
	Dispatch DispatchConfig
	Conn     ConnConfig
	Logic    LogicConfig
	App AppConfig
}

type AppConfig struct {
	GetServiceRetrySeconds int
}

// DispatchConfig dispatch的配置
type DispatchConfig struct {
	AutoReConnection         bool // 断开连接时是否自动重新连接
	ConnectionIntervalSecond int  // 断连时重新连接的间隔时间
	ConnectionRetryCount     int  // 重试次数

	GetPingTime  time.Duration
	GetPongTime  time.Duration
	SendPingTime time.Duration
}

// LogicConfig logic层的配置
type LogicConfig struct {
	RequestTimeOut int // 逻辑层中一个请求发来，等待响应时的超时时间
	Tracing TracingConfig
}

type TracingConfig struct {
	Enable bool
	Sugar bool
}

// ConnConfig 连接的config
type ConnConfig struct {
	ReadTimeout int // 读一个数据包时的超时时间，用在一个数据包未读完时
}

func init() {
	viper.SetDefault("DISPATCH_CONFIG", "{}")
}

// 加载远程配置
func remoteConfig() envConfig {
	// 1. 启动begonia同样需要配置 => 1.1 使用本地默认配置启动 1.2 拉取远程配置

	// 配置中心
	// - 提供配置的增删改查
	// - 可以拉取配置 key(serviceName) => value  namespace(db_config) => value
	// DispatchConfig, ConnConfig, AppConfig
	// 1. 拉取默认配置 namespace => value    func FetchConfigByNamespace()
	// 2. 拉取服务特定配置 namespace, serviceName => value    func FetchConfigByNamespaceAndServiceName()
	// 3. 拉取服务配置 serviceName => value   func FetchConfigByServiceName()

	/*
	"envConfig":{},
	"dbConfig":{},
	"businessConfig":{},
	"redisConfig":{},

	1. 拉取
	2. 同步到绑定的结构体
	3. 如果有新的结构体被绑定，将数据同步上去

	1. 添加服务端对客户端主动推送的功能 / 服务中心对节点推送的功能
	2. 与框架整合
	config.Watch("redisConfig",func())


	1. Framework 初始化阶段
	config.Bind("envConfig",&C)

	2. 业务
	config.Bind("redisConfig",&redis.Config)
	 */
	return envConfig{}
}

// setupConfig 加载本地配置
func setupConfig() envConfig {
	return envConfig{}
}

func defaultConfig() envConfig {
	return envConfig{
		App: AppConfig{
			GetServiceRetrySeconds: 1,
		},
		Conn: ConnConfig{
			ReadTimeout: 10,
		},
		Dispatch: DispatchConfig{
			AutoReConnection:         true,
			ConnectionIntervalSecond: 2,
			ConnectionRetryCount:     1000,
			GetPingTime:              20 * time.Second,
			GetPongTime:              10 * time.Second,
			SendPingTime:             10 * time.Second,
		},
		Logic: LogicConfig{
			RequestTimeOut: 10,
			Tracing: TracingConfig{
				Enable: true,
				Sugar:  false,
			},
		}}
}
