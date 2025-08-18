package redisclient

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client redis.UniversalClient
}

func NewRedisClientFromURL(rawURL string) (*RedisClient, error) {
	logger := log.New(os.Stdout, "redis-client", 1)

	cleanURL := cleanRedisURL(rawURL)

	parsed, err := parseURL(cleanURL)

	if err != nil {
		logger.Print("Parsed error to : %w", err)
		return nil, fmt.Errorf("parsed error to : %w", err)
	}

	password := extractPassword(parsed)
	addrs := extractAddrs(parsed)

	switch parsed.Scheme {
	case "redis":
		return createRedisClient(addrs, password, 0), nil

	case "redis+sentinel", "sentinel":
		logger.Println("connect into redis sentinel mode")
		return createSentinelClient(rawURL, parsed, password), nil

	case "redis+cluster", "cluster":
		logger.Println("connect into redis cluster")
		return createClusterClient(rawURL, password), nil

	default:
		if len(addrs) == 0 {
			return nil, fmt.Errorf("invalid redis URL: %s", rawURL)
		}

		return createRedisClient(addrs, password, 0), nil
	}
}

func (r *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func cleanRedisURL(rawURL string) string {
	if strings.HasPrefix(rawURL, "http://") {
		return strings.TrimPrefix(rawURL, "http://")
	} else if strings.HasPrefix(rawURL, "https://") {
		return strings.TrimPrefix(rawURL, "https://")
	}
	return rawURL
}

func parseURL(urlStr string) (*url.URL, error) {
	return url.Parse(urlStr)
}

func extractPassword(parsed *url.URL) string {
	if parsed.User != nil {
		password, _ := parsed.User.Password()
		return password
	}
	return ""
}

func extractAddrs(parsed *url.URL) []string {
	return strings.Split(parsed.Host, ",")
}

func createRedisClient(addrs []string, password string, db int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:         addrs[0],
		Password:     password,
		DB:           db,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	return &RedisClient{client: client}
}

func createSentinelClient(rawURL string, parsed *url.URL, password string) *RedisClient {
	hosts := strings.Split(strings.Split(strings.Split(rawURL, "//")[1], "/")[0], ",")

	path := strings.TrimPrefix(parsed.Path, "/")

	masterName := path

	if strings.Contains(path, "service_name:") {
		masterName = strings.Split(path, "service_name:")[1]
	}

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: hosts,
		Password:      password,
		DB:            0,
		PoolSize:      20,
		MinIdleConns:  5,
		DialTimeout:   2 * time.Second,
		ReadTimeout:   1 * time.Second,
		WriteTimeout:  1 * time.Second,
	})

	return &RedisClient{client: client}
}

func createClusterClient(rawURL string, password string) *RedisClient {
	hosts := strings.Split(strings.Split(strings.Split(rawURL, "//")[1], "/")[0], ",")

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        hosts,
		Password:     password,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	return &RedisClient{client: client}
}
