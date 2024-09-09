package handlers

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"SimpleApp/global"
)

func PostEchoHandler(c echo.Context) error {
	time.Sleep(global.State.DelayInMS)
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	log.Println("Received Message: ", string(bodyBytes))
	return c.String(http.StatusOK, string(bodyBytes))
}

func GetHealthCheckHandler(c echo.Context) error {
	time.Sleep(global.State.DelayInMS)
	return c.String(http.StatusOK, "HEALTHY")
}
