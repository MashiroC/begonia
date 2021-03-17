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
}

// DispatchConfig dispatch的配置
type DispatchConfig struct {
	AutoReConnection           bool // 断开连接时是否自动重新连接
	ReConnectionIntervalSecond int  // 断连时重新连接的间隔时间
	ReConnectionRetryCount     int  // 重试次数

	GetPingTime  time.Duration
	GetPongTime  time.Duration
	SendPingTime time.Duration
}

// LogicConfig logic层的配置
type LogicConfig struct {
	RequestTimeOut int // 逻辑层中一个请求发来，等待响应时的超时时间
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
	return envConfig{}
}

// setupConfig 加载本地配置
func setupConfig() envConfig {
	return envConfig{}
}

func defaultConfig() envConfig {
	return envConfig{
		Conn: ConnConfig{
			ReadTimeout: 10,
		},
		Dispatch: DispatchConfig{
			AutoReConnection:           true,
			ReConnectionIntervalSecond: 2,
			ReConnectionRetryCount:     5,
			GetPingTime:                20 * time.Second,
			GetPongTime:                10 * time.Second,
			SendPingTime:               10 * time.Second,
		},
		Logic: LogicConfig{
			RequestTimeOut: 10,
		}}
}
