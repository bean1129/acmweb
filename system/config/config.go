package config

const (
	// EnvName ConfigEnv 配置环境
	EnvName = "CONFIG"
	// FilePath ConfigFile 配置文件
	FilePath = "conf/config.yaml"
)

// ResourceConfig 全局配置结构体
type ResourceConfig struct {
	Mysql       MySQL       `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis       Redis       `mapstructure:"Redis" json:"Redis" yaml:"Redis"`
	Attachment  Attachment  `mapstructure:"attachment" json:"attachment" yaml:"attachment"`
	Application Application `mapstructure:"application" json:"application" yaml:"application"`
}

var CONFIG ResourceConfig
