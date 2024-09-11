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

	httpClient := http.Client{
		Timeout: global.NODE_MAX_TIMEOUT,
	}

	for {
		selectedNode := global.State.GetNextHealthyNode()
		if selectedNode == "" {
			break
		}

		body, err := passRequest(
			&httpClient,
			selectedNode+c.Request().URL.Path,
			c.Request().Header.Get("Content-Type"),
			bytes.NewReader(bodyBytes),
		)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		log.Println("Active Node: [" + strings.Join(global.State.HealthyNodes, ",") + "]")

		return c.String(http.StatusOK, string(body))
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "Instance not available")
}

var passRequest = func(httpClient *http.Client, url string, header string, body *bytes.Reader) (response []byte, err error) {
	res, _ := httpClient.Post(
		url,
		header,
		body,
	)

	defer res.Body.Close()
	response, err = io.ReadAll(res.Body)

	return response, err
}
