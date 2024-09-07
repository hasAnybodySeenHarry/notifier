package redis

import "github.com/redis/go-redis/v9"

type Clients struct {
	Users UserModel
}

func New(c *redis.Client) *Clients {
	return &Clients{
		Users: UserModel{client: c},
	}
}
