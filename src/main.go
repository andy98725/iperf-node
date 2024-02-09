package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	if useQuickstart(&s) {
		return
	}

	if err = s.connect(); err != nil {
		e.Logger.Fatal(err)
		panic(err)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "healthy"})
	})
	e.POST("/start", func(c echo.Context) error {
		body := make(map[string]interface{})
		if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
		}

		key, ok := body["key"].(string)
		if !ok {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "missing key"})
		}
		if !s.validate(key) {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "invalid key"})
		}

		testIdRaw, ok := body["testId"]
		if !ok {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "missing testId"})
		}
		testId, err := strconv.Atoi(fmt.Sprintf("%v", testIdRaw))
		if err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "testId must be integer"})
		}

		mode, ok := body["mode"].(string)
		if !ok {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "missing mode"})
		}

		if mode == "client" {
			port, ok := body["port"].(string)
			if !ok {
				return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "missing port"})
			}
			serverAddress, ok := body["serverAddress"].(string)
			if !ok {
				return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "missing serverAddress"})
			}

			if err := s.runIperfClient(testId, serverAddress, port, s.completeClientTest); err != nil {
				return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
			}
			return c.String(http.StatusOK, "iPerf server started successfully")
		} else {
			if err := s.runIperfServer(testId); err != nil {
				return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
			}
			return c.String(http.StatusOK, "iPerf server started successfully")
		}
	})
	e.POST("/finish", func(c echo.Context) error {
		if err := s.closeIperf(); err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
		}

		if err := s.completeServerTest(); err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
		}

		return c.String(http.StatusOK, "iPerf server closed successfully")
	})
	e.POST("/close", func(c echo.Context) error {
		if err := s.closeIperf(); err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
		}
		// Refresh connection status
		if err := s.connect(); err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
		}

		return c.String(http.StatusOK, "iPerf closed successfully")
	})

	e.Logger.Fatal(e.Start(":" + s.config.hostPort))
}
