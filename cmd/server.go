package cmd

import (
	"fmt"

	"github.com/chatmcp/mcprouter/router"
	"github.com/chatmcp/mcprouter/service/sse"
	"github.com/spf13/cobra"
)

var port int

func startServer() {
	s := sse.NewSSEServer()

	s.Route(router.Route)

	s.Start(port)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start sse server",
	Long:  `start sse server as MCP Server Proxy`,
	Run: func(cmd *cobra.Command, args []string) {
		if port == 0 {
			fmt.Println("port is required")
			return
		}
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVarP(&port, "port", "p", 8025, "port to run the server on")
}
