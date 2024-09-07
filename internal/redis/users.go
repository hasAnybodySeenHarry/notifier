package redis

import "github.com/redis/go-redis/v9"

type UserModel struct {
	client *redis.Client
}
