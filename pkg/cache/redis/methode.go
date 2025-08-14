package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"todo-app/internal/repository/mysql/models"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RediseCLientDB struct {
	Db *redis.Client
}

func NewRedisDB(db *redis.Client) *RediseCLientDB {
	return &RediseCLientDB{Db: db}
}

type Recording struct {
	ID    int
	List  models.TodoList
	Items models.TodoItems
}

func (r *RediseCLientDB) Create(record *Recording) error {

	if r == nil || r.Db == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	var (
		key        string
		recordJSON []byte
		err        error
	)

	if record.List.Title != "" {
		key = "list"
		recordJSON, err = json.Marshal(record.List)
		if err != nil {
			return err
		}
	} else if record.Items.Title != "" {
		key = "item"
		recordJSON, err = json.Marshal(record.Items)
		if err != nil {
			return err
		}
	}

	_, err = r.Get(record.ID, key)
	if err != nil {
		return err
	}

	const timeOfLife = 30 * time.Second
	if err := r.Db.Set(context.Background(), key+"_"+fmt.Sprint(record.ID), recordJSON, timeOfLife).Err(); err != nil {
		return err
	}

	return nil

}

func (r *RediseCLientDB) Get(key int, typeKey string) (*Recording, error) {

	if r == nil || r.Db == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}

	idStr := strconv.Itoa(key)
	if typeKey == "list" || typeKey == "item" {
		idStr = typeKey + "_" + idStr

		logrus.Debug("Get redis key: ", idStr)

	} else {
		return nil, fmt.Errorf("error invalid key redis")
	}

	ctx := context.Background()
	val, err := r.Db.Get(ctx, idStr).Result()
	if err == redis.Nil {
		return &Recording{}, nil
	} else if err != nil {
		return &Recording{}, err
	}

	var recording *Recording
	err = json.Unmarshal([]byte(val), &recording)
	if err != nil {
		return &Recording{}, err
	}

	return recording, nil

}
