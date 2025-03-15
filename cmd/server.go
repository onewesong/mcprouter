package cmd

import (
	"fmt"

	"github.com/chatmcp/mcprouter/router"
	"github.com/chatmcp/mcprouter/service/sse"
	"github.com/chatmcp/mcprouter/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

func startServer(port int) {
	s := sse.NewSSEServer()

	s.Route(router.Route)
	s.Start(port)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start sse server",
	Long:  `start sse proxy server`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := util.InitConfigWithFile(configFile); err != nil {
			fmt.Printf("init config failed with file: %s, %v\n", configFile, err)
			return
		}

		port := viper.GetInt("server.port")
		if port == 0 {
			port = 8025
		}

		fmt.Printf("starting server on port: %d\n", port)
		startServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&configFile, "config", "c", ".env.toml", "config file (default is .env.toml)")
}
