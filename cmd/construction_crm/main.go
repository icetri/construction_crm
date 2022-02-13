package main

import (
	"flag"
	"github.com/construction_crm/internal/clients/postgres"
	"github.com/construction_crm/internal/construction_crm/server"
	"github.com/construction_crm/internal/construction_crm/server/handlers"
	"github.com/construction_crm/internal/construction_crm/service"
	"github.com/construction_crm/internal/construction_crm/service/firebase"
	"github.com/construction_crm/internal/construction_crm/service/mail"
	"github.com/construction_crm/internal/construction_crm/service/minio"
	"github.com/construction_crm/internal/construction_crm/types/config"
	"github.com/construction_crm/pkg/logger"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	configPath := new(string)
	flag.StringVar(configPath, "config-path", "config/config-local.yaml", "specify path to yaml")
	flag.Parse()

	configFile, err := os.Open(*configPath)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with os.Open config"))
	}

	cfg := config.Config{}
	if err := yaml.NewDecoder(configFile).Decode(&cfg); err != nil {
		logger.LogFatal(errors.Wrap(err, "err with Decode config"))
	}

	if err = logger.NewLogger(cfg.Telegram); err != nil {
		logger.LogFatal(err)
	}

	//для запуска необходимо прописать в ./config/config-local.yaml,
	//если это тестовый запуск то необходимо закомментировать иницаилизацию NewMinio
	fs, err := minio.NewMinio(cfg.FileStorage)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewMinio"))
	}

	//для запуска необходимо прописать в ./firebase_example.json,
	//если это тестовый запуск то необходимо закомментировать иницаилизацию NewFireBase
	fb, err := firebase.NewFireBase(cfg.FireBase)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewFireBase"))
	}

	db, err := postgres.NewPostgres(cfg.PostgresDsn)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewPostgres"))
	}
	defer db.Close()

	//для запуска необходимо прописать в ./config,
	//если это тестовый запуск то необходимо закомментировать иницаилизацию NewMail
	email, err := mail.NewMail(cfg.Email, db)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewMail"))
	}

	s := service.NewService(db, &cfg, fs, fb, email, cfg.SMS)

	endpoints := handlers.NewHandlers(s, &cfg)

	server.StartServer(endpoints, cfg.ServerPort)
}
