package main

import (
	"errors"
	"os"
	"os/exec"
	"strconv"

	"github.com/labstack/echo"
)

type server struct {
	config struct {
		hostPort   string
		iperfPort  string
		serverAddr string
		id         int
		key        string
		hash       string
	}

	log     echo.Logger
	process *exec.Cmd
	testId  int
}

func initServer(log echo.Logger) (server, error) {
	s := server{log: log}

	s.config.hostPort = os.Getenv("PORT")
	if s.config.hostPort == "" {
		s.config.hostPort = "8080"
	}
	s.config.iperfPort = os.Getenv("IPERF_PORT")
	if s.config.iperfPort == "" {
		s.config.iperfPort = "5001"
	}

	idStr := os.Getenv("ID")
	if idStr == "" {
		return s, errors.New("env ID is required")
	}
	i, err := strconv.Atoi(idStr)
	if err != nil {
		return s, errors.New("env ID must be integer")
	}
	s.config.id = i

	s.config.serverAddr = os.Getenv("ENDPOINT")
	if s.config.serverAddr == "" {
		return s, errors.New("env ENDPOINT is required")
	}
	s.config.key = os.Getenv("ENDPOINT_KEY")
	if s.config.key == "" {
		return s, errors.New("env ENDPOINT_KEY is required")
	}
	s.config.hash = os.Getenv("HASH")
	if s.config.hash == "" {
		return s, errors.New("env HASH is required")
	}

	return s, s.connect()
}

func (s *server) validate(key string) bool {
	// TODO validate
	return true
}
