package service

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Service struct {
	Config Config
	Logger Logger
}

func New() (s Service, err error) {
	// Get config file path.
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}

	// Load config file and instanciate logger.
	s.Config, err = LoadConfig(configPath)
	if err != nil {
		return s, fmt.Errorf("load config: %w", err)
	}
	s.Logger, err = NewLogger(s.Config.Log)
	if err != nil {
		return s, fmt.Errorf("init logger: %w", err)
	}

	s.Logger.Log("Loaded config.\n")
	s.Logger.Log("Logger ready.\n")

	return s, nil
}

func (s Service) Shutdown() (err error) {
	s.Logger.Log("Shutting down...\n")

	// Sync and close log file if regular file.
	if f, ok := s.Logger.Writer.(*os.File); ok && f != os.Stdout && f != os.Stderr {
		s.Logger.Log("Synchronizing log file...\n")
		err = f.Sync()
		if err != nil {
			s.Logger.Log(err.Error() + "\n")
			return err
		}
		s.Logger.Log("Closing log file...\n")
		err = f.Close()
		if err != nil {
			s.Logger.Log(err.Error() + "\n")
			return err
		}
	} else {
		s.Logger.Log("Shutdown completed.\n")
	}

	return nil
}

func (s Service) Run() {
	s.Logger.Log("Running...\n")
	s.Logger.Log("Current time: " + time.Now().Format(time.RFC3339) + "\n")
	s.Logger.Log("PID: " + strconv.Itoa(os.Getpid()) + "\n")
	s.Logger.Log("Environment: " + s.Config.Env + "\n")
	s.Logger.Log("Log file: " + s.Config.Log + "\n")

	for t := range time.Tick(5 * time.Second) {
		s.Logger.Log(t.Format(time.RFC3339) + "\n")
	}
}
