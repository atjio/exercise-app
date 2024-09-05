package main 

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main () {
	portNumber := flag.String("port", "1234", "Service port")
	flag.Parse()

	port := fmt.Sprintf(":%s", *portNumber)

	e := echo.New()
	e.POST("/echo", func (c echo.Context) error {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.String(http.StatusOK, string(bodyBytes))
	})

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "HEALTHY")
	})


	e.Logger.Fatal(e.Start(port))
}