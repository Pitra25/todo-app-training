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

type RedisCLientDB struct {
	Db *redis.Client
}

func NewRedisDB(db *redis.Client) *RedisCLientDB {
	return &RedisCLientDB{Db: db}
}

type Recording struct {
	ID       int
	List     models.TodoList
	Items    models.TodoItems
	CodeUser models.UsersCode
}

func (r *RedisCLientDB) Create(record *Recording) error {

	if r == nil || r.Db == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	var (
		key        keyType
		recordJSON []byte
		err        error
		timeOfLife = 30 * time.Second
	)

	switch {
	case record.List.Title != "":
		{
			key = List
			recordJSON, err = json.Marshal(record.List)
			if err != nil {
				return err
			}
		}
	case record.Items.Title != "":
		{
			key = Item
			recordJSON, err = json.Marshal(record.Items)
			if err != nil {
				return err
			}
		}
	case record.CodeUser.Code != "":
		{
			key = Code_user
			timeOfLife = 10 * time.Minute
			recordJSON, err = json.Marshal(record.CodeUser)
			if err != nil {
				return err
			}
		}
	}

	keyRdb := keyFormation(record.ID, key)
	if keyRdb != "" {
		return fmt.Errorf("error create key recording redis")
	}

	_, err = r.Get(record.ID, key)
	if err != nil {
		return err
	}

	err = r.Db.SetNX(
		context.Background(),
		keyRdb,
		recordJSON,
		timeOfLife,
	).Err()
	if err != nil {
		return err
	}

	return nil

}

func (r *RedisCLientDB) Get(key int, typeKey keyType) (*Recording, error) {

	if r == nil || r.Db == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}

	keyRdb := keyFormation(key, typeKey)
	if keyRdb != "" {
		return nil, fmt.Errorf("error create key recording redis")
	}

	ctx := context.Background()
	val, err := r.Db.Get(ctx, keyRdb).Result()
	if err == redis.Nil {
		logrus.Debug("Error redis: ", err.Error())
		return &Recording{}, nil
	} else if err != nil {
		logrus.Debug("Error get redis: ", err.Error())
		return &Recording{}, err
	}

	var recording *Recording
	err = json.Unmarshal([]byte(val), &recording)
	if err != nil {
		return &Recording{}, err
	}

	logrus.Debug("check redis: ", recording.CodeUser.Code)

	return recording, nil

}

func (r *RedisCLientDB) DeleteRecord(id int, typeKey keyType) error {
	keyRdb := keyFormation(id, typeKey)
	if keyRdb != "" {
		return fmt.Errorf("error create key recording redis")
	}

	ctx := context.Background()
	_, err := r.Db.Del(ctx, keyRdb).Result()
	return err
}

type keyType int

const (
	List keyType = iota
	Item
	Code_user
)

func keyFormation(id int, typeKey keyType) string {
	idUserStr := strconv.Itoa(id)
	key := ""

	switch typeKey {
	case List:
		{
			key = "list_" + idUserStr
		}
	case Item:
		{
			key = "itm_" + idUserStr
		}
	case Code_user:
		{
			key = "code_user_" + idUserStr
		}
	default:
		{
			logrus.Error("error invalid key redis")
			return ""
		}
	}

	return key
}
