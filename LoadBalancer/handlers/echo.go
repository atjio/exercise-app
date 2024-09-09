package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"LoadBalancer/global"
)

func PostEchoHandler(c echo.Context) error {
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if (len(global.State.HealthyNodes) == 0) {
		return echo.NewHTTPError(http.StatusInternalServerError, "Instance not available")
	}

	for i := 0; i < len(global.State.HealthyNodes); i++ {
		selectedNode := global.State.HealthyNodes[i]
		res, err := http.Post(
			selectedNode + c.Request().URL.Path, 
			c.Request().Header.Get("Content-Type"), 
    		bytes.NewReader(bodyBytes),
		)
		
		if (err != nil) {
			continue
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		global.State.HealthyNodes = append(global.State.HealthyNodes[i+1:], global.State.HealthyNodes[:i+1]...)
		log.Println("Active Node: [" + strings.Join(global.State.HealthyNodes, ",") + "]")

		return c.String(http.StatusOK, string(body))
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "Instance not available")
}