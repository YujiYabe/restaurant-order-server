package redis

import (
	"net"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/YujiYabe/restaurant-order-server/backend/internal/2_adapter/gateway"
)

type (
	// Redis は、Redisへの接続を保持します。
	Redis struct {
		Conn *redis.Client
	}
)

// NewToRedis は、必要な依存関係と初期値を設定したインスタンスを生成します。
func NewToRedis() (
	toRedis gateway.ToRedis,
) {
	conn := redis.NewClient(&redis.Options{
		Addr:     redisAddress(),
		Password: "", // パスワードなし
		DB:       0,  // 既定のDBを使用
	})

	redis := new(Redis)
	redis.Conn = conn

	toRedis = redis
	return
}

func redisAddress() string {
	address := strings.TrimSpace(os.Getenv("REDIS_ADDR"))
	if address != "" {
		return address
	}
	port := strings.TrimSpace(os.Getenv("REDIS_PORT"))
	if port == "" {
		port = "6379"
	}
	return net.JoinHostPort("127.0.0.1", port)
}

// init は、パッケージの初期状態を構成します。
func init() {
}
