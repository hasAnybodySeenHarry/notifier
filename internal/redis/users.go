package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	client *redis.Client
}

type NotiBroadcast struct {
	ID       primitive.ObjectID
	Type     string
	ForUsers []int64
}

func (m *UserModel) Publish(n *NotiBroadcast) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	js, err := json.Marshal(n)
	if err != nil {
		return err
	}

	_, err = m.client.Publish(ctx, "noti-broadcast", js).Result()
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Subscribe(channel string) *redis.PubSub {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.client.Subscribe(ctx, channel)
}
