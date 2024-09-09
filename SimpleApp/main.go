package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v4"

	"SimpleApp/global"
	"SimpleApp/handlers"
)

func main() {
	portNumber := flag.String("port", "1234", "Service port")
	flag.Parse()
	global.PORT = fmt.Sprintf(":%s", *portNumber)

	e := echo.New()
	e.POST("/echo", handlers.PostEchoHandler)
	e.GET("/healthcheck", handlers.GetHealthCheckHandler)
	e.GET("/debug", handlers.GetDebugHandler)
	e.POST("/simulateDelay", handlers.PostSimulateDelayHandler)

	go handlers.RegisterService(global.LOAD_BALANCER_URL)

	e.Logger.Fatal(e.Start(global.PORT))
}
