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
	log    echo.Logger
	config struct {
		name       string
		hostPort   string
		iperfPort  string
		serverAddr string
	}

	process   *exec.Cmd
	outStr    string
	outBuffer *bytes.Buffer
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

func (s *server) connect() error {
	body := struct {
		Name string
	}{Name: s.config.name}

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

	s.log.Debug("response recieved:")
	s.log.Debug(resp.Status)
	return nil
}

func (s *server) runIperf() error {
	if s.process != nil {
		return errors.New("iPerf is already running")
	}

	s.process = exec.Command("iperf", "-s -p "+s.config.iperfPort)
	s.outStr = ""
	s.outBuffer = new(bytes.Buffer)
	s.process.Stdout = s.outBuffer

	s.process.Start()

	return nil
}
func (s *server) getIperfState() (string, error) {
	if s.process == nil {
		return "", errors.New("iPerf is not running yet")
	}

	return s.outBuffer.String(), nil
}
