package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/chatmcp/mcprouter/util"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const (
	proxyInfoKey = "pi_%s"
)

func getRedisClient() (*redis.Client, error) {
	cli := util.GetRedisClient(viper.GetString("app.cache_name"))
	if cli == nil {
		return nil, fmt.Errorf("redis client not found")
	}

	return cli, nil
}

func getRedisContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// StoreProxyInfo stores the proxy info to redis
func StoreProxyInfo(sessionID string, proxyInfo *ProxyInfo) error {
	ctx, cancel := getRedisContext()
	defer cancel()

	cli, err := getRedisClient()
	if err != nil {
		return err
	}

	b, err := json.Marshal(proxyInfo)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(proxyInfoKey, sessionID)
	expires := 30 * 24 * time.Hour // 30 days
	cli.Set(ctx, cacheKey, b, expires)

	return nil
}

// GetProxyInfo gets the proxy info from redis
func GetProxyInfo(sessionID string) (*ProxyInfo, error) {
	ctx, cancel := getRedisContext()
	defer cancel()

	cli, err := getRedisClient()
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf(proxyInfoKey, sessionID)

	b, err := cli.Get(ctx, cacheKey).Bytes()

	if err != nil || b == nil {
		return nil, errors.New("cache not found")
	}

	proxyInfo := &ProxyInfo{}
	if err := json.Unmarshal(b, proxyInfo); err != nil {
		return nil, err
	}

	return proxyInfo, nil
}

func DeleteProxyInfo(sessionID string) error {
	ctx, cancel := getRedisContext()
	defer cancel()

	cli, err := getRedisClient()
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf(proxyInfoKey, sessionID)
	cli.Del(ctx, cacheKey)

	return nil
}
