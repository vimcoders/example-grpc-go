package balance

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func newUniversalClient(_ context.Context, addrs []string) redis.UniversalClient {
	opts := redis.UniversalOptions{
		Addrs: addrs,
		// 每个节点的连接池大小
		PoolSize: 1,
		// 集群重定向最大次数（自动发现时可能需要）
		MaxRedirects: 3,
		// 节点失败标记时间（自动避开故障节点）
		FailingTimeoutSeconds: 15,
		// 按延迟路由（可选，提升性能）
		RouteByLatency: true,
	}
	return redis.NewUniversalClient(&opts)
}
