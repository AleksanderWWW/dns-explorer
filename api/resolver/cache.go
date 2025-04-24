package resolver

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachedDNSResponse struct {
	Records []string // IPs, CNAMEs, etc.
}

func CacheDNSResponse(ctx context.Context, rdb *redis.Client, key string, response CachedDNSResponse, ttl uint32) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}

func GetFromCache(ctx context.Context, rdb *redis.Client, key string) (*CachedDNSResponse, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, err
	}

	var res CachedDNSResponse
	if err := json.Unmarshal([]byte(val), &res); err != nil {
		return nil, err
	}
	return &res, nil
}
