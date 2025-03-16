package mcpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// GetCommand returns the command for the given key
func GetCommand(key string) string {
	command := viper.GetString(fmt.Sprintf("mcp_server_commands.%s", key))
	if command == "" {
		return getRemoteCommand(key)
	}

	return command
}

// getRemoteCommand returns the command for the given key from the remote API
func getRemoteCommand(key string) string {
	apiUrl := viper.GetString("remote_apis.get_server_command")

	params := map[string]string{
		"server_key": key,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return ""
	}

	response, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return ""
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	data := gjson.ParseBytes(body)
	command := data.Get("data.server_command").String()

	return command
}
