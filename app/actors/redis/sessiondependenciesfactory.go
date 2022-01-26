package redis

import (
	"context"

	"github.com/garyburd/redigo/redis"
)

func SessionDependenciesFactory(command string, pool redis.Pool) (func(args ...interface{}) (interface{}, int), redis.Conn) {
	ctx := context.Background()
	client := pool.Get()
	// Build dependency inversion
	return Build(client, ctx, command), client
}
