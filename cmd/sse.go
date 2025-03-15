/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

func startSSEServer() {
	// Create MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
	)

	// Add tool
	tool := mcp.NewTool("ping",
		mcp.WithDescription("Ping the server"),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("Message to ping the server"),
		),
	)

	// Add tool handler
	s.AddTool(tool, pingHandler)

	// Start the stdio server
	sseServer := server.NewSSEServer(s, server.WithBaseURL("http://127.0.0.1:8027"))
	log.Printf("SSE server listening on :8027")
	if err := sseServer.Start(":8027"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func pingHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	message, ok := request.Params.Arguments["message"].(string)
	if !ok {
		return mcp.NewToolResultError("message must be a string"), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Pong: %s!", message)), nil
}

// sseCmd represents the sse command
var sseCmd = &cobra.Command{
	Use:   "sse",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sse called")
		startSSEServer()
	},
}

func init() {
	rootCmd.AddCommand(sseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
