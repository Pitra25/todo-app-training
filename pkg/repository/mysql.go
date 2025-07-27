package repository

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const ( // table DB
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	TodoItemssTable = "todo_items"
	listsItemsTable = "lists_items"
)

type ConfigMySql struct {
	Host      string
	Port      string
	Username  string
	Password  string
	DBName    string
	SSLMode   string
	ParseTime string
}

func NewMySqlDB(cfg ConfigMySql) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.ParseTime)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
