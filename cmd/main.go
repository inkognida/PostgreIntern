package main

import (
	"github.com/sirupsen/logrus"
	"postgreintern/internal/repository"
	services "postgreintern/internal/service"
)

func main() {
	logger := logrus.New()

	// repo - уровень работы с бд
	repo := repository.NewRepo(logger)
	// service - сервис для обработки изменений в директориях
	service := services.NewService(logger, repo)

	// установка конфига с директориями и командами
	err := service.SetupConfigWatcher()
	if err != nil {
		logger.Fatalln("error during config setup/watcher creation", err)
	}

	// основной процесс сервиса
	service.ExecuteProcess()
}

