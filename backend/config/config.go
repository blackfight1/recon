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

	// 允许从环境变量读取配置
	viper.AutomaticEnv()
	viper.SetEnvPrefix("") // 不使用前缀

	// 绑定环境变量到配置项
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	AppConfig = &Config{}
	return viper.Unmarshal(AppConfig)
}
