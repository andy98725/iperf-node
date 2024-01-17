package main

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	s := initState()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/start", func(c echo.Context) error {
		err := s.runIperf()
		if err != nil {
			return c.JSON(http.StatusOK, struct{ Error string }{Error: err.Error()})
		}

		return c.String(http.StatusOK, "iPerf started successfully on port "+s.config.iperfPort)
	})
	e.GET("/status", func(c echo.Context) error {
		state, err := s.getIperfState()
		if err != nil {
			return c.JSON(http.StatusOK, struct{ Error string }{Error: err.Error()})
		}

		return c.String(http.StatusOK, state)

	})

	e.Logger.Fatal(e.Start(":" + s.config.serverPort))
}

type serverState struct {
	config struct {
		serverPort string
		iperfPort  string
	}

	process   *exec.Cmd
	outStr    string
	outBuffer *bytes.Buffer
}

func (s *serverState) runIperf() error {
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
func (s *serverState) getIperfState() (string, error) {
	if s.process == nil {
		return "", errors.New("iPerf is not running yet")
	}

	return s.outBuffer.String(), nil
}

func initState() serverState {
	s := serverState{}

	s.config.serverPort = os.Getenv("PORT")
	if s.config.serverPort == "" {
		s.config.serverPort = "8080"
	}
	s.config.iperfPort = os.Getenv("IPERF_PORT")
	if s.config.iperfPort == "" {
		s.config.iperfPort = "5001"
	}

	return s
}
