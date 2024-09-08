package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func postEchoHandler(c echo.Context) error {
	time.Sleep(delayInMS)
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	log.Println("Received Message: ", string(bodyBytes))
	return c.String(http.StatusOK, string(bodyBytes))
}

func getHealthCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "HEALTHY")
}

func getDebugHandler(c echo.Context) error { 
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/register", nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	req.Header.Set("X-Client-Port", port)

	res, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, string(body))
}

func postSimulateDelayHandler(c echo.Context) error {
	d := c.FormValue("delay")
	delay, err := strconv.Atoi(d)
	if err != nil {
		delayInMS = 0 * time.Millisecond
	} else {
		delayInMS = time.Duration(delay) * time.Millisecond
	}

	return c.String(http.StatusOK, "Delay = " + string(delayInMS))
}