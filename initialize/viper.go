package initialize

import (
	"fmt"
	"makedb/global"

	"github.com/spf13/viper"
)

func InitViper(configFile string) {
	v := viper.New()
	v.SetConfigFile(configFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error read config file: %s", err))
	}
	if err = v.Unmarshal(&global.MAKEDB_CONFIG); err != nil {
		panic(fmt.Errorf("fatal error Unmarshal config file: %s", err))
	}
}
