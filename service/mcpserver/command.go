package mcpserver

import (
	"fmt"

	"github.com/spf13/viper"
)

// GetCommand returns the command for the given key
func GetCommand(key string) string {
	command := viper.GetString(fmt.Sprintf("mcp_server_commands.%s", key))
	if command == "" {
		return ""
	}

	return command
}
