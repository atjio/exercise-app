package main 

import (
	"flag"
	"fmt"
	"time"
	
	"github.com/labstack/echo/v4"
)

var (
	port string
	loadBalancerUrl string = "http://127.0.0.1:8080"
	delayInMS time.Duration = 0 * time.Millisecond
)

func main () {
	portNumber := flag.String("port", "1234", "Service port")
	flag.Parse()
	port = fmt.Sprintf(":%s", *portNumber)

	e := echo.New()
	e.POST("/echo", postEchoHandler)
	e.GET("/healthcheck", getHealthCheckHandler)
	e.GET("/debug", getDebugHandler)
	e.POST("/simulateDelay", postSimulateDelayHandler)
	
	go func() {
		registerService(loadBalancerUrl + "/register")
	}()

	e.Logger.Fatal(e.Start(port))
}
