package app

import (
	"os"
	"todo-app/internal/repository/mysql"
	storage "todo-app/pkg/cache/redis"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func bdInit() *sqlx.DB {
	logrus.Print("Database initialization started")

	db, err := mysql.NewMySqlDB(
		&mysql.ConfigMySql{
			Host:      viper.GetString("db.host"),
			Port:      viper.GetString("db.port"),
			Username:  viper.GetString("db.username"),
			DBName:    viper.GetString("db.dbname"),
			SSLMode:   viper.GetString("db.sslmode"),
			ParseTime: viper.GetString("db.parsetime"),
			Password:  os.Getenv("DB_PASSWORD"),
		})
	if err != nil {
		logrus.Fatalln("failed to initialize db: ", err.Error())
		return nil
	}

	logrus.Print("Database initialization completed")

	return db
}

func redisInit() *redis.Client {
	logrus.Print("Redis initialization started")

	rdb, err := storage.New(storage.ConfigRedis{
		Addr:        viper.GetString("redis.host") + viper.GetString("redis.port"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		User:        viper.GetString("redis.user"),
		DB:          viper.GetInt("redis.db"),
		MaxRetries:  viper.GetInt("redis.maxretries"),
		DialTimeout: viper.GetDuration("redis.dialtimeout"),
		Timeout:     viper.GetDuration("redis.timeout"),
	})
	if err != nil {
		logrus.Fatalln("failed to initialize redis: ", err.Error())
		return nil
	}

	logrus.Print("Redis initialization completed")

	return rdb
}
