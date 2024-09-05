package main 

import (
	"flag"
	"fmt"
	
	"github.com/labstack/echo/v4"
)

var port string

func main () {
	portNumber := flag.String("port", "1234", "Service port")
	flag.Parse()
	port = fmt.Sprintf(":%s", *portNumber)

	e := echo.New()
	e.POST("/echo", postEchoHandler)
	e.GET("/healthcheck", getHealthCheckHandler)
	
	e.GET("/debug", getDebugHandler)

	e.Logger.Fatal(e.Start(port))
}
