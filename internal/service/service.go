package services

import (
	"context"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"postgreintern/internal/model"
	"postgreintern/internal/repository"
	"postgreintern/internal/utils"
	"strconv"
	"sync"
	"syscall"
)

const (
	RunCommandError = "run error"
	WatcherError = "watcher error"
)

type Service struct {
	// атрибуты сервиса
	logger *logrus.Logger
	config model.Config

	// работа с бд
	repo repository.Repository

	// синхронизация сервиса
	stopCh chan struct{}
	wg *sync.WaitGroup
}

func NewService(logger *logrus.Logger, repo repository.Repository) *Service {
	return &Service{
		logger:  logger,
		repo:    repo,
		wg:      &sync.WaitGroup{},
		stopCh: make(chan struct{}),
	}
}

func (s *Service) SetupConfigWatcher() error {
	viper.SetConfigFile("../config.yaml")

	// чтение файла конфигурации
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// сохранение конфигурации в структуре Config
	err = viper.Unmarshal(&s.config)
	if err != nil {
		return err
	}

	return nil
}



func (s *Service) ExecuteProcess() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGKILL)

	for _, dir := range s.config.Dirs {
		// инициализируем wg для каждой директории
		s.wg.Add(1)
		// отслеживание изменений
		dirCopy := dir
		go func(repo repository.Repository, dir *model.PathConfig) {
			defer s.wg.Done()

			if utils.CurrentDir(dirCopy.Path, s.logger) {
				s.logger.Infoln("You cant use current directory")
				return
			}

			err := s.Execute(repo, dir)
			if err != nil {
				s.logger.Infoln(err)
				return
			}
		}(s.repo, &dirCopy)
	}

	<-signalCh
	s.logger.Infoln("Ctrl+C signal received. Stopping goroutines...")
	go func() {
		s.wg.Wait()
		close(s.stopCh)
	}()

	s.logger.Infoln("All commands finished")
}

func (s *Service) Execute(repo repository.Repository, dir *model.PathConfig) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(dir.Path)
	if err != nil {
		s.logger.Fatalln(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New(WatcherError)
			} else { // какое-то изменение
				// запись изменения и получение id для выполнения команд
				id, err := repo.SaveEvent(context.Background(), model.FileEvent{
					EventType: utils.GetOpString(event.Op),
					Path:      dir.Path,
					FileName:  utils.GetLastTextBefore(event.Name, "/"),
				})
				if err != nil {
					s.logger.Fatalln("failed to save file changes to db", err)
				}
				// запуск команд и логирование (внутри RunCommands)
				s.logger.Infoln("Event:", event.Op, event.Name)
				err = s.RunCommands(dir, id)
				if err != nil {
					return err
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
		}
	}
}

func (s *Service) RunCommands(dir *model.PathConfig, id int) error {
	workDir := dir.Path

	unique := strconv.Itoa(id)

	file, err := os.Create("../logs/log"+unique)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	for i, _ := range dir.Commands {
		fullCmd := dir.Commands[i].Cmd
		cmd := exec.Command("bash", "-c", fullCmd)
		cmd.Dir = workDir
		cmd.Stdout = file

		done := make(chan struct{})
		err := cmd.Start()
		if err != nil {
			log.Printf("Error executing command '%s': %v\n", cmd, err)
			return errors.New("error executing")
		}

		go func() {
			select {
			case <-s.stopCh:
				// сигнал для остановки горутины, убиваем процесс
				if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
					return
				}
			case <-done:
				// горутина выполнена
			}
		}()

		err = cmd.Wait()
		if err != nil {
			s.logger.Infoln("FAILED: specification of error for debug process:", err)
			close(done)
			// возращаем общую ошибку выполнения команд для остановки выполнения
			return errors.New(RunCommandError)
		} else {
			execCmd, args := utils.GetCmdArgs(fullCmd)
			err = s.repo.SaveCommandExecution(context.Background(), model.CmdEvent{
				Cmd:  execCmd,
				Args: args,
			}, id)
			if err != nil {
				s.logger.Fatalln("failed to save file changes to db", err)
			} else {
				s.logger.Infoln("SUCCESS:", dir.Commands[i].Cmd)
			}
			close(done)
		}
	}
	return nil
}