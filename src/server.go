package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo"
)

type server struct {
	config struct {
		name       string
		hostPort   string
		iperfPort  string
		serverAddr string
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
	s.config.serverAddr = os.Getenv("ENDPOINT")
	if s.config.serverAddr == "" {
		return s, errors.New("env ENDPOINT is required")
	}
	s.config.name = os.Getenv("NAME")
	if s.config.name == "" {
		return s, errors.New("env NAME is required")
	}

	return s, s.connect()
}

func (s *server) validate(key string) bool {
	// TODO validate
	return true
}

func (s *server) connect() error {
	body := struct {
		Name       string
		ServerPort string
		iPerfPort  string
	}{
		Name:       s.config.name,
		ServerPort: s.config.hostPort,
		iPerfPort:  s.config.iperfPort,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.config.serverAddr+"/api/nodes/connect", &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resp := &struct{ Status string }{}
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return err
	}

	if res.StatusCode != 200 {
		s.log.Debug("response recieved:")
		s.log.Debug(resp.Status)
	}

	return nil
}

func (s *server) completeTest(results *string) error {
	var buf bytes.Buffer

	if results != nil {
		body := struct {
			Results string
			TestId  int
		}{Results: *results, TestId: s.testId}
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	} else {
		body := struct {
			TestId int
		}{TestId: s.testId}
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest("POST", s.config.serverAddr+"/api/nodes/complete", &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resp := &struct{ Status string }{}
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return err
	}

	if res.StatusCode != 200 {
		s.log.Debug("response recieved:")
		s.log.Debug(resp.Status)
	}

	return nil
}
