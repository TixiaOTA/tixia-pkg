package cache

import (
	"time"

	"github.com/TixiaOTA/tixia-pkg/redis/client"
)

type CacheConfig struct {
	BaseRedisConfig client.BaseRedisConfig
	DefaultTTL      time.Duration
	MaxTTL          time.Duration
	KeyPrefix       string
	EnableMetrics   bool
}

type CacheOptions struct {
	TTL         time.Duration
	KeyPrefix   string
	Compression bool
}
