package base

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RdbBaseInterface interface {
	VerifyConnection() (bool, error)
}

// VerifyConnection 验证redis的连接
func VerifyConnection(rDb *redis.Client) (bool, error) {
	_, err := rDb.Ping(context.TODO()).Result()
	return err == nil, err
}
