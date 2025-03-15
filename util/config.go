package util

import (
	"github.com/spf13/viper"
)

// InitConfigWithFile will read config from file
func InitConfigWithFile(filename string) error {
	viper.SetConfigFile(filename)

	return viper.ReadInConfig()
}
