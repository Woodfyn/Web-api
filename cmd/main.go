package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Woodfyn/Web-api/internal/config"
	"github.com/Woodfyn/Web-api/internal/handler/rest"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
	"github.com/Woodfyn/Web-api/internal/service"
	"github.com/Woodfyn/Web-api/pkg/database"
	"github.com/Woodfyn/Web-api/pkg/server"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {

	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		logrus.Fatalf("config is not initialised: %s", err.Error())
	}

	db, err := database.NewPostgesDB(database.ConnInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		logrus.Fatalf("config was not transferred to the db: %s", err.Error())
	}

	repos := psql.NewRepository(db)
	services := service.NewService(repos)
	handlers := rest.NewHandler(services)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg.Server.Port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("the port is not specified in the configuration: %s", err.Error())
		}
	}()

	logrus.Info("SERVER STARTED")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("SERVER STOPPED")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("errer occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("errer occured on db connection close %s", err.Error())
	}
}
