package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/AleksanderWWW/dns-explorer/resolver"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	domain := "example.com"

	// Check cache first
	cached, err := resolver.GetFromCache(ctx, rdb, domain)
	if err != nil {
		panic(err)
	}

	if cached != nil {
		fmt.Println("Cache hit:", cached.Records)
		return
	}

	// Resolve and cache
	res, ttl, err := resolver.ResolveDNS(domain)
	if err != nil {
		panic(err)
	}

	fmt.Println("Resolved:", res.Records, "TTL:", ttl)
	if err := resolver.CacheDNSResponse(ctx, rdb, domain, res, ttl); err != nil {
		panic(err)
	}
}
