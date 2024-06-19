package config

// Application 系统配置
type Application struct {
	Name       string `mapstructure:"name" json:"name" yaml:"name"`
	Version    string `mapstructure:"version" json:"version" yaml:"version"`
	Model      int    `mapstructure:"model" json:"model" yaml:"model"`
	Addr       string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Port       string `mapstructure:"port" json:"port" yaml:"port"`
	FilePort   int    `mapstructure:"fileport" json:"fileport" yaml:"fileport"`
	Debug      bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	Image      string `mapstructure:"image" json:"image" yaml:"image"`
	Uploads    string `mapstructure:"uploads" json:"uploads" yaml:"uploads"`
	SecretKey  string `mapstructure:"secret_key" jsoin:"secret_key" yaml:"secret_key"`
	ServerId   int    `mapstructure:"serv_id" json:"serv_id" yaml:"serv_id"`
	ExpireTime int    `mapstructure:"expire_time" json:"expire_time" yaml:"expire_time"`
	LogLevel   string `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	CheckTime  int    `mapstructure:"check_time" json:"check_time" yaml:"check_time"`
}
