package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func init() {
	var filename string
	if configEnv := os.Getenv(EnvName); configEnv == "" {
		filename = FilePath
		fmt.Printf("您正在使用config的默认值,config的路径为%v\n", FilePath)
	} else {
		filename = configEnv
		fmt.Printf("您正在使用CONFIG环境变量,config的路径为%v\n", filename)
	}

	v := viper.New()
	v.SetConfigFile(filename)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&CONFIG); err != nil {
			panic(err)
		}
	})

	if err := v.Unmarshal(&CONFIG); err != nil {
		panic(err)
	}
}
