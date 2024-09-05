package main

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"log"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
)

func getRegisterHandler(c echo.Context) error {
	clientPort := c.Request().Header.Get("X-Client-Port")

	if (clientPort == "") {
		return echo.NewHTTPError(http.StatusInternalServerError, "Client Port is required")
	}

	clientIP, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if (err != nil) {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	url := "http://" + clientIP + clientPort

	if (!slices.Contains(state.healthyNodes, url)) {
		res, err := http.Get(url + "/healthcheck")
		if (err != nil) {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if (res.StatusCode == http.StatusOK) {
			state.healthyNodes = append(state.healthyNodes, url)
		
		}
	}

	log.Println("Active Node: [" + strings.Join(state.healthyNodes, ",") + "]")
	

	return c.String(http.StatusOK, url)
}

func postEchoHandler(c echo.Context) error {
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var bodyJson map[string]interface{}

	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if (len(state.healthyNodes) == 0) {
		return echo.NewHTTPError(http.StatusInternalServerError, "Instance not available")
	}

	for i := 0; i < len(state.healthyNodes); len++ {
		res, err := http.Post(state.healthyNodes[i] + "/echo", "application/json", bodyJson);
		
		if (err != nil) {
			continue
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		state.healthyNodes = append(state.healthyNodes[i:], state.healthyNodes[:i])
		log.Println("Active Node: [" + strings.Join(state.healthyNodes, ",") + "]")

		return c.String(http.StatusOK, string(body))
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "Instance not available")
}