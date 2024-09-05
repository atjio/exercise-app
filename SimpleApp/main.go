package main 

import (
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
)

func main () {
	portNumber := flag.String("port", "1234", "Service port")
	flag.Parse()
	port := fmt.Sprintf(":%s", *portNumber)

	e := echo.New()
	e.POST("/echo", PostEchoHandler)
	e.GET("/healthcheck", GetHealthCheckHandler)

	e.Logger.Fatal(e.Start(port))
}
