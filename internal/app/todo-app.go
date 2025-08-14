package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"todo-app"
	v1 "todo-app/internal/handler/http/v1"
	"todo-app/internal/repository"
	"todo-app/internal/service"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func TodoApp() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db := bdInit()
	rClient := redisInit()
	
	repos := repository.NewRepository(db, rClient)
	service := service.NewService(repos)
	handler := v1.NewHandler(service)

	srv := startApp(handler)
	closeApp(srv, db)

}

func startApp(h *v1.Heandler) *todo.Server {
	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("server.port"), h.InitRoutes()); err != nil {
			logrus.Fatalf("error run http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")
	logrus.Printf(
		"URL: http://%s:%s/%s",
		viper.GetString("server.host"),
		viper.GetString("server.port"),
		viper.GetString("server.swagger_url"),
	)

	return srv
}

func closeApp(srv *todo.Server, db *sqlx.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
