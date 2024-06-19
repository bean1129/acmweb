package config

// Redis 结构体
type Redis struct {
	Network   string `mapstructure:"Network" json:"Network" yaml:"Network"`
	Addr      string `mapstructure:"Addr" json:"Addr" yaml:"Addr"`
	Timeout   int    `mapstructure:"Timeout" json:"Timeout" yaml:"Timeout"`
	MaxActive int    `mapstructure:"MaxActive" json:"MaxActive" yaml:"MaxActive"`
	Password  string `mapstructure:"Password" json:"Password" yaml:"Password"`
	Database  string `mapstructure:"Database" json:"Database" yaml:"Database"`
	Prefix    string `mapstructure:"Prefix" json:"Prefix" yaml:"Prefix"`
	Delim     string `mapstructure:"Delim" json:"Delim" yaml:"Delim"`
}
