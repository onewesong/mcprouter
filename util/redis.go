package util

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RedisConfig is the config for the redis
type RedisConfig struct {
	Host     string
	Port     int64
	Password string
	DB       int
	PoolSize int
}

var redisCH = make(chan map[string]*redis.Client)

// InitRedisWithName initializes the redis connection
func InitRedisWithName(name string) error {
	var conf RedisConfig
	sub := viper.Sub("cache." + name)
	if sub == nil {
		return fmt.Errorf("invalid redis config under %s", name)
	}
	if err := sub.Unmarshal(&conf); err != nil {
		return err
	}

	if conf.Port == 0 {
		conf.Port = 6379
	}

	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	}

	if conf.PoolSize > 0 {
		opts.PoolSize = conf.PoolSize
	}

	cli := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return err
	}

	// set cli
	redisCH <- map[string]*redis.Client{name: cli}

	return nil
}

// GetClient gets the redis connection client
func GetRedisClient(name string) *redis.Client {
	climap := <-redisCH
	if cli, ok := climap[name]; ok {
		return cli
	}

	return nil
}

func cliPool() {
	var clis = make(map[string]*redis.Client)
	for {
		select {
		case climap := <-redisCH:
			for name, cli := range climap {
				clis[name] = cli
			}
		case redisCH <- clis:
		}
	}
}

func init() {
	go cliPool()
}
