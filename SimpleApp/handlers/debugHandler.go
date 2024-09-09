package handlers

import (
	"io"
	"net/http"
	"strconv"
	"time"
	
	"github.com/labstack/echo/v4"

	"SimpleApp/global"
)

func GetDebugHandler(c echo.Context) error { 
	client := &http.Client{}

	req, err := http.NewRequest("GET", global.REGISTER_URL, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	req.Header.Set("X-Client-Port", global.PORT)

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

func PostSimulateDelayHandler(c echo.Context) error {
	d := c.FormValue("delay")
	delay, err := strconv.Atoi(d)
	if err != nil {
		global.State.DelayInMS = 0 * time.Millisecond
	} else {
		global.State.DelayInMS = time.Duration(delay) * time.Millisecond
	}

	return c.String(http.StatusOK, "Delay = " + global.State.DelayInMS.String())
}