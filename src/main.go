package main

import (
	"encoding/json"
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

			if err := s.runIperfClient(serverAddress, port); err != nil {
				return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
			}
			return c.String(http.StatusOK, "iPerf server started successfully")
		} else {
			if err := s.runIperfServer(); err != nil {
				return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: err.Error()})
			}
			return c.String(http.StatusOK, "iPerf server started successfully")

		}
	})

	e.Logger.Fatal(e.Start(":" + s.config.hostPort))
}
