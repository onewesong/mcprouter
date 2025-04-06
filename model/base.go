package model

import (
	"github.com/chatmcp/mcprouter/util"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func db() *gorm.DB {
	name := viper.GetString("app.web_db_name")

	return util.GetClient(name)
}

func adb() *gorm.DB {
	name := viper.GetString("app.api_db_name")

	return util.GetClient(name)
}
