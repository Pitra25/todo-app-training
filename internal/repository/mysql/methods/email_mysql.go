package methods

import (
	"fmt"
	"time"
	"todo-app/internal/repository/mysql/models"
	storage "todo-app/pkg/cache/redis"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type EmailMySql struct {
	db  *sqlx.DB
	rdb *storage.RedisCLientDB
}

func NewEmailMySql(db *sqlx.DB, rdb *redis.Client) *EmailMySql {
	return &EmailMySql{
		db:  db,
		rdb: &storage.RedisCLientDB{Db: rdb},
	}
}

func (e *EmailMySql) SaveCodeUser(code string, userId int) error {

	// Save to Redis cache
	if err := e.rdb.Create(&storage.Recording{
		ID: userId,
		CodeUser: models.UsersCode{
			UserId:    userId,
			Code:      code,
			ExpiresAt: "10m",
		},
	}); err != nil {
		logrus.Error("err creatr record redis.", err.Error())
		return nil
	}

	createCodeQuery := fmt.Sprintf(
		"INSERT INTO %s (user_id, code, expires_at, date_of_creation) VALUES (?, ?, ?, ?)",
		models.UserCodeEmailTable,
	)

	timeOfLife := time.Now().Add(10 * time.Minute).Format(time.DateTime)
	timeCreate := time.Now().Format(time.DateTime)

	_, err := e.db.Exec(createCodeQuery, userId, code, timeOfLife, timeCreate)
	if err != nil {
		return err
	}

	return nil
}

type ResponseCode struct {
	Id   int
	Code string
}

func (e *EmailMySql) GetCodeUser(userId int) (ResponseCode, error) {

	// Check Redis cache first
	storageRecords, err := e.rdb.Get(userId, storage.Code_user)
	if err != nil {
		logrus.Error("GetCodeUSer code user.", err.Error())
	} else if storageRecords.CodeUser.Code != "" && storageRecords != nil {
		logrus.Debug("GetCodeUSer code user from redis cache. id:", userId)
		return ResponseCode{
			Code: storageRecords.CodeUser.Code,
		}, nil
	}

	logrus.Debug("not redis")

	// Get from db mysql
	getCodeQuery := fmt.Sprintf(
		"SELECT id code FROM %s WHERE user_id = ? AND expires_at > NOW()",
		models.UserCodeEmailTable,
	)

	var rCode ResponseCode
	err = e.db.Get(&rCode, getCodeQuery, userId)
	if err != nil {
		logrus.Debug("Error: ", err.Error())
		return ResponseCode{}, err
	}

	logrus.Debug("check select DB")

	return ResponseCode{
		Id:   rCode.Id,
		Code: rCode.Code,
	}, nil
}

func (e *EmailMySql) UpdateStatusUser(userId int) error {

	tx, err := e.db.Begin()
	if err != nil {
		return err
	}

	changeStatusQuery := fmt.Sprintf(
		"UPDATE %s SET status = ? WHERE user_id = ?",
		models.UsersTable,
	)

	_, err = tx.Exec(changeStatusQuery, "confirmed", userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (e *EmailMySql) DeleteRecord(id, userId int) error {
	deleteQuery := fmt.Sprintf(
		"DELETE us FROM %s us WHERE us.id = ? AND us.user_id = ?",
		models.UserCodeEmailTable,
	)

	_, err := e.db.Exec(deleteQuery, id, userId)

	return err
}
