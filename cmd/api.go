package cmd

import (
	"fmt"
	"log"

	"github.com/chatmcp/mcprouter/router"
	"github.com/chatmcp/mcprouter/service/api"
	"github.com/chatmcp/mcprouter/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiConfigFile string

// startAPIServer starts the api server
func startAPIServer(port int) {
	s := api.NewAPIServer()

	s.Route(router.APIRoute)
	s.Start(port)
}

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start api server",
	Long:  `start api server`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := util.InitConfigWithFile(apiConfigFile); err != nil {
			fmt.Printf("init config failed with file: %s, %v\n", apiConfigFile, err)
			return
		}

		log.Println("config initialized")

		if viper.GetBool("app.use_db") && viper.GetString("app.db_name") != "" {
			if err := util.InitDBWithName(viper.GetString("app.db_name")); err != nil {
				fmt.Printf("init db failed with name: %s, %v\n", viper.GetString("app.db_name"), err)
				return
			}
			log.Println("db initialized")
		}

		if viper.GetBool("app.use_cache") && viper.GetString("app.cache_name") == "redis" {
			if err := util.InitRedisWithName(viper.GetString("app.cache_name")); err != nil {
				fmt.Printf("init redis failed with name: %s, %v\n", viper.GetString("app.cache_name"), err)
				return
			}
			log.Println("redis initialized")
		}

		port := viper.GetInt("api_server.port")
		if port == 0 {
			port = 8027
		}

		startAPIServer(port)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	apiCmd.Flags().StringVarP(&apiConfigFile, "config", "c", ".env.toml", "config file (default is .env.toml)")
}
