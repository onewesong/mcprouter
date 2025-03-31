package mcpserver

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// GetServerCommand returns the command for the given key
func GetServerCommand(key string) string {
	config := GetServerConfig(key)
	if config == nil {
		return ""
	}

	return config.Command
}

// GetServerConfig returns the config for the given key
func GetServerConfig(key string) *ServerConfig {
	config := &ServerConfig{}
	err := viper.UnmarshalKey(fmt.Sprintf("mcp_servers.%s", key), config)

	if config.Command == "" {
		fmt.Printf("get local config failed: %v, try to get remote config\n", err)

		config, err = getRemoteServerConfig(key)
		if err != nil {
			fmt.Printf("get remote config failed: %v\n", err)
			return nil
		}
	}

	return config
}

// getRemoteServerConfig returns the config for the given key from the remote API
func getRemoteServerConfig(key string) (*ServerConfig, error) {
	apiUrl := viper.GetString("remote_apis.get_server_config")

	params := map[string]string{
		"server_key": key,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	data := gjson.ParseBytes(body)
	if data.Get("code").Int() != 0 {
		return nil, fmt.Errorf("get remote config failed: %s", data.Get("message").String())
	}

	config := &ServerConfig{
		ServerUUID:   data.Get("data.server_uuid").String(),
		ServerName:   data.Get("data.server_name").String(),
		ServerKey:    data.Get("data.server_key").String(),
		Command:      data.Get("data.command").String(),
		CommandHash:  data.Get("data.command_hash").String(),
		ShareProcess: data.Get("data.share_process").Bool(),
	}

	if config.CommandHash == "" {
		config.CommandHash = fmt.Sprintf("%x", md5.Sum([]byte(config.Command)))
	}

	return config, nil
}
