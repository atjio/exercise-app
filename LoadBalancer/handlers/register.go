package handlers

import (
	"net"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"

	"LoadBalancer/global"
)

func GetRegisterHandler(c echo.Context) error {
	registeredNodes := global.State.GetAllNodes()
	clientPort := c.Request().Header.Get("X-Client-Port")

	if clientPort == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Client Port is required")
	}

	clientIP, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	url := global.HTTP_PROTOCOL + clientIP + clientPort

	if !slices.Contains(registeredNodes, url) && triggerHealthcheck(url) {
		global.State.AddNode(url)
	}

	return c.String(http.StatusOK, url)
}
