package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ejuju/my-website/pkg/httpmux"
)

type Service struct {
	Config        Config
	Logger        Logger
	HTTPEndpoints httpmux.Endpoints
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

	s.Logger.Log("Config and logger ready.")
	return s, nil
}

func (s Service) Shutdown() (err error) {
	// Sync and close log file if regular file.
	if f, ok := s.Logger.Writer.(*os.File); ok && f != os.Stdout && f != os.Stderr {
		s.Logger.Log("Closing logfile...")
		err = f.Sync()
		if err != nil {
			s.Logger.Log(err.Error())
		}

		s.Logger.Log("Closing log file...")
		err = f.Close()
		if err != nil {
			return fmt.Errorf("close log file: %w", err)
		}
		return nil
	}

	s.Logger.Log("Goodbye!")
	return nil
}

func (s *Service) Run() {
	s.Logger.Log("At:          " + time.Now().Format(time.RFC3339))
	s.Logger.Log("PID:         " + strconv.Itoa(os.Getpid()))
	s.Logger.Log("Log file:    " + s.Config.Log)
	s.Logger.Log("Environment: " + s.Config.Env)
	s.Logger.Log("HTTP Port:   " + strconv.Itoa(s.Config.HTTPPort))
	s.Logger.Log("CTCP Port:   " + strconv.Itoa(s.Config.CTCPPort))

	go s.runHTTPServer()
}

func (s Service) Panic(err error) {
	s.Logger.Log("Panic: " + err.Error())
	s.Logger.Log("Shutting down and panicking...")
	shutdownErr := s.Shutdown()
	if shutdownErr != nil {
		s.Logger.Log("Shutdown failed: " + shutdownErr.Error())
	}
	panic(err)
}
