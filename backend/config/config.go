package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	Scanner      ScannerConfig      `mapstructure:"scanner"`
	Notification NotificationConfig `mapstructure:"notification"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type ScannerConfig struct {
	Interval    int    `mapstructure:"interval"`
	DataDir     string `mapstructure:"data_dir"`
	Concurrency int    `mapstructure:"concurrency"`
}

type NotificationConfig struct {
	Wecom    WecomConfig    `mapstructure:"wecom"`
	Dingtalk DingtalkConfig `mapstructure:"dingtalk"`
}

type WecomConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Webhook string `mapstructure:"webhook"`
}

type DingtalkConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Webhook string `mapstructure:"webhook"`
	Secret  string `mapstructure:"secret"`
}

var AppConfig *Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	AppConfig = &Config{}
	return viper.Unmarshal(AppConfig)
}
