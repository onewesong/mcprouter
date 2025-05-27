package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chatmcp/mcprouter/router"
	"github.com/chatmcp/mcprouter/service/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiConfigFile string

// startAPIServer starts the api server with graceful shutdown support
func startAPIServer(port int) {
	s := api.NewAPIServer()
	s.Route(router.APIRoute)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("API server starting on port %d", port)

	// Start server with context for graceful shutdown
	if err := s.StartWithContext(ctx, port); err != nil {
		log.Printf("API server failed to start: %v", err)
	}

	log.Println("API server stopped gracefully")
}

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start api server",
	Long:  `start api server`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Init(); err != nil {
			log.Printf("init failed with error: %v", err)
			return
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
