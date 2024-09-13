package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Clients struct {
	Users UserModel
}

func New(c *redis.Client) *Clients {
	return &Clients{
		Users: UserModel{client: c},
	}
}
