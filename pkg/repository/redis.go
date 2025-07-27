package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"todo-app/types"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type ConfigRedis struct {
	Addr        string        `yml:"addr"`
	Password    string        `yml:"password"`
	User        string        `yml:"user"`
	DB          int           `yml:"db"`
	MaxRetries  int           `yml:"max_retries"`
	DialTimeout time.Duration `yml:"dial_timeout"`
	Timeout     time.Duration `yml:"timeout"`
}

var cfg = ConfigRedis{
	Addr:        viper.GetString("redis.host") + viper.GetString("redis.host"),
	Password:    os.Getenv("REDIS_PASSWORD"),
	User:        viper.GetString("redis.user"),
	DB:          viper.GetInt("redis.db"),
	MaxRetries:  viper.GetInt("redis.maxretries"),
	DialTimeout: viper.GetDuration("redis.dialtimeout"),
	Timeout:     viper.GetDuration("redis.timeout"),
}

func NewClientRedis(ctx *context.Context) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(*ctx).Err(); err != nil {
		return nil, err
	}

	return db, nil
}

type RecordingStruct struct {
	ID    int
	list  types.TodoList
	items types.TodoItems
}

func Create(record RecordingStruct) error {
	var (
		key        string
		recordJSON []byte
		err        error
	)

	if record.list.Title != "" {
		key = "list"
		recordJSON, err = json.Marshal(record.list)
		if err != nil {
			return err
		}
	} else if record.items.Title != "" {
		key = "item"
		recordJSON, err = json.Marshal(record.items)
		if err != nil {
			return err
		}
	}

	_, err = Get(record.ID, key)
	if err != nil {
		return err
	}

	ctx := context.Background()
	dbR, err := NewClientRedis(&ctx)
	if err != nil {
		return err
	}

	const timeOfLife = 30 * time.Second
	if err := dbR.Set(context.Background(), key+"_"+fmt.Sprint(record.ID), recordJSON, timeOfLife).Err(); err != nil {
		return err
	}
	return nil
}

func Get(key int, typeKey string) (*RecordingStruct, error) {
	ctx := context.Background()
	rdb, err := NewClientRedis(&ctx)
	if err != nil {
		return &RecordingStruct{}, err
	}

	idStr := strconv.Itoa(key)
	if typeKey == "list" || typeKey == "item" {
		idStr = typeKey + "_" + idStr
	} else {
		return nil, fmt.Errorf("error invalid key redis")
	}

	val, err := rdb.Get(ctx, idStr).Result()
	if err == redis.Nil {
		return &RecordingStruct{}, nil
	} else if err != nil {
		return &RecordingStruct{}, err
	}

	var recording *RecordingStruct
	err = json.Unmarshal([]byte(val), &recording)
	if err != nil {
		return &RecordingStruct{}, err
	}

	return recording, nil
}
