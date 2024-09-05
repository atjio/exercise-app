package main

import (
	"io"
	"net/http"
	"github.com/labstack/echo/v4"
)

func PostEchoHandler(c echo.Context) error {
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.String(http.StatusOK, string(bodyBytes))
}

func GetHealthCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "HEALTHY")
}