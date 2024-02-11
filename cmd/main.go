package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Woodfyn/Web-api/internal/config"
	"github.com/Woodfyn/Web-api/internal/handler/rest"
	"github.com/Woodfyn/Web-api/internal/repository/psql"
	"github.com/Woodfyn/Web-api/internal/service"
	"github.com/Woodfyn/Web-api/pkg/auth"
	"github.com/Woodfyn/Web-api/pkg/database"
	"github.com/Woodfyn/Web-api/pkg/hash"
	"github.com/Woodfyn/Web-api/pkg/server"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

// @title GameList API
// @version 1.0
// @description API Server for GameList Application

// @host localhost:8000
// @BasePath /api

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
	CONFIG_ENV  = ".main"
	LOG_FILE    = "logfile.log"
	LOG_FOLDER  = "logs"
)

func init() {
	logsDir := fmt.Sprintf("./%s", LOG_FOLDER)
	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		logrus.Fatal(err)
	}

	logFile := filepath.Join(logsDir, LOG_FILE)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetOutput(file)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("recovered from panic: %v", err)
		}

		logrus.StandardLogger().Out.(*os.File).Close()
	}()
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE, CONFIG_ENV)
	logrus.Info(cfg)
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

	hasher := hash.NewSHA1Hasher(cfg.Hash.Salt)
	tokenManager, err := auth.NewManager(cfg.Auth.Secret)
	if err != nil {
		logrus.Fatalf("No init token manager: %s", err.Error())
	}

	repos := psql.NewRepositories(db)

	deps := service.Deps{
		Repos:           repos,
		Hasher:          hasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.JWT.AccessTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTTL,
	}

	service := service.NewServices(deps)

	handlers := rest.NewHandler(service, tokenManager)

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

	if err := srv.Shutdown(context.TODO()); err != nil {
		logrus.Errorf("errer occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("errer occured on db connection close %s", err.Error())
	}
}
