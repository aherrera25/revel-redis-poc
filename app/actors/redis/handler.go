package redis

import (
	"context"

	"github.com/revel/revel"
)

type SessionProvider interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
}

// factory method providing access to redis features in a SOLID way
func Build(client SessionProvider, ctx context.Context, command string) func(args ...interface{}) (interface{}, int) {
	// provide a behavior function
	return func(args ...interface{}) (interface{}, int) {
		currMessage, err := client.Do(command, args...)
		if err == nil {
			return currMessage, 200
		}
		revel.AppLog.Error(err.Error())
		return err.Error(), 500
	}
}
