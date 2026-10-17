package redis

import "testing"

func TestRedisAddress(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		t.Setenv("REDIS_ADDR", "")
		t.Setenv("REDIS_PORT", "")
		if got := redisAddress(); got != "127.0.0.1:6379" {
			t.Fatalf("redisAddress() = %q, want %q", got, "127.0.0.1:6379")
		}
	})

	t.Run("custom port", func(t *testing.T) {
		t.Setenv("REDIS_ADDR", "")
		t.Setenv("REDIS_PORT", "16379")
		if got := redisAddress(); got != "127.0.0.1:16379" {
			t.Fatalf("redisAddress() = %q, want %q", got, "127.0.0.1:16379")
		}
	})

	t.Run("address takes precedence", func(t *testing.T) {
		t.Setenv("REDIS_ADDR", "redis.internal:6380")
		t.Setenv("REDIS_PORT", "16379")
		if got := redisAddress(); got != "redis.internal:6380" {
			t.Fatalf("redisAddress() = %q, want %q", got, "redis.internal:6380")
		}
	})
}
