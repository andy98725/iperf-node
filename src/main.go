package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)

	s, err := initServer(e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
		panic(err)
	}

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

	e.Logger.Fatal(e.Start(":" + s.config.hostPort))
}
