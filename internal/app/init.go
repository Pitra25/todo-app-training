package app

import (
	"net/smtp"
	"os"
	"todo-app/internal/repository/mysql"
	"todo-app/internal/repository/postgres"
	storage "todo-app/pkg/cache/redis"
	"todo-app/pkg/email"

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

func bdMysqlInit() *sqlx.DB {
	logrus.Print("Database initialization started")

	db, err := mysql.New(
		&mysql.ConfigMySql{
			Host:      viper.GetString("db.mysql.host"),
			Port:      viper.GetString("db.mysql.port"),
			Username:  viper.GetString("db.mysql.username"),
			DBName:    viper.GetString("db.mysql.dbname"),
			SSLMode:   viper.GetString("db.mysql.sslmode"),
			ParseTime: viper.GetString("db.mysql.parsetime"),
			Password:  os.Getenv("DB_PASSWORD"),
		})
	if err != nil {
		logrus.Fatalln("failed to initialize db: ", err.Error())
		return nil
	}

	logrus.Print("Database initialization completed")

	return db
}

func bdPostgresInit() *sqlx.DB {
	logrus.Print("Database initialization started")

	db, err := postgres.New(
		&postgres.ConfigPostgres{
			Host:      viper.GetString("db.mysql.host"),
			Port:      viper.GetString("db.mysql.port"),
			Username:  viper.GetString("db.mysql.username"),
			DBName:    viper.GetString("db.mysql.dbname"),
			SSLMode:   viper.GetString("db.mysql.sslmode"),
			ParseTime: viper.GetString("db.mysql.parsetime"),
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

func initConfigSMTP() *smtp.Client {
	logrus.Print("SMTP configuration initialization started")

	smtpCLient := email.New(email.ConfigSMTP{
		Host:     viper.GetString("smtp.emailClient.yandex.host"),
		Port:     viper.GetString("smtp.emailClient.yandex.port_smtp_starttls"),
		Username: viper.GetString("smtp.emailClient.yandex.username"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     viper.GetString("smtp.from"),
	})

	if smtpCLient == nil {
		logrus.Fatalln("failed to initialize SMTP configuration")
		return nil
	}

	logrus.Print("SMTP configuration initialization completed")
	return smtpCLient
}
