package config

// MySQL 数据库结构体
type MySQL struct {
	Host        string `mapstructure:"host" json:"host" yaml:"host"`
	Port        int    `mapstructure:"port" json:"port" yaml:"port"`
	Database    string `mapstructure:"database" json:"database" yaml:"database"`
	Username    string `mapstructure:"username" json:"username" yaml:"username"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	Charset     string `mapstructure:"charset" json:"charset" yaml:"charset"`
	Debug       bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	LogLevel    string `mapstructure:"logLevel" json:"logLevel" yaml:"logLevel"`
	MaxOpenCons int    `mapstructure:"maxOpenCons" json:"maxOpenCons" yaml:"maxOpenCons"`
	MaxIdleCons int    `mapstructure:"maxIdleCons" json:"maxIdleCons" yaml:"maxIdleCons"`
	Timeout     string `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	ParseTime   bool   `mapstructure:"parseTime" json:"parseTime" yaml:"parseTime"`
}
