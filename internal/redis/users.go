package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

type State int

const (
	JOINED State = iota
	CONNECTED
	DISCONNECTED
	REMOVED
)

type User struct {
	ID    string
	State State
}

func (s State) String() string {
	switch s {
	case JOINED:
		return "JOINED"
	case DISCONNECTED:
		return "DISCONNECTED"
	case REMOVED:
		return "REMOVED"
	default:
		return "UNKNOWN"
	}
}

func (m *UserModel) AddUserState(userID int64, state State) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userIDStr := strconv.FormatInt(userID, 10)

	val, err := m.client.HGet(ctx, "users", userIDStr).Result()
	if err != nil {
		return err
	}

	var userState State
	switch val {
	case "JOINED":
		userState = JOINED
	case "DISCONNECTED":
		userState = DISCONNECTED
	case "REMOVED":
		userState = REMOVED
	case "CONNECTED":
		userState = CONNECTED
	default:
		return fmt.Errorf("unknown state: %v", val)
	}

	if userState != JOINED {
		return ErrUserAlreadyExists
	}

	return m.client.HSet(ctx, "users", userIDStr, state.String()).Err()
}

func (m *UserModel) InitUserState(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userIDStr := strconv.FormatInt(userID, 10)

	val, err := m.client.HGet(ctx, "users", userIDStr).Result()
	if err == redis.Nil {
		return m.client.HSet(ctx, "users", userIDStr, JOINED.String()).Err()
	} else if err != nil {
		return err
	}

	var userState State
	switch val {
	case "JOINED":
		userState = JOINED
	case "DISCONNECTED":
		userState = DISCONNECTED
	case "REMOVED":
		userState = REMOVED
	case "CONNECTED":
		userState = CONNECTED
	default:
		return fmt.Errorf("unknown state: %v", val)
	}

	if userState == JOINED || userState == CONNECTED {
		return ErrUserAlreadyExists
	}

	// If the user is in another state (DISCONNECTED or REMOVED), set to JOINED
	return m.client.HSet(ctx, "users", userIDStr, JOINED.String()).Err()
}

func (m *UserModel) DeleteUser(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userIDStr := strconv.FormatInt(userID, 10)

	_, err := m.client.HGet(ctx, "users", userIDStr).Result()
	if err == redis.Nil {
		return fmt.Errorf("user not found")
	} else if err != nil {
		return err
	}

	return m.client.HDel(ctx, "users", userIDStr).Err()
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
