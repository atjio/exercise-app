package main

import (
	"fmt"
	"flag"
	"time"
	"github.com/labstack/echo/v4"
)

type loadBalancerStatus struct {
	healthyNodes []string 
	unhealthyNodes []string 
}

// type node struct {
// 	url string
// }

// func healthcheck()

var (
	state loadBalancerStatus
	healthcheckUrl string = "/healthcheck"
	healtcheckDelayInMS time.Duration = 1000 * time.Millisecond
	nodeMaxTimeout time.Duration = 2000 * time.Millisecond
)
func main() {
	portNumber := flag.String("port", "8080", "LoadBalancer Port")
	flag.Parse()
	port := fmt.Sprintf(":%s", *portNumber)

	e := echo.New()

	e.GET("/register", getRegisterHandler)
	e.POST("/*", postEchoHandler)

	go func() {
		for {
			healthyNodes, unhealthyNodes := healthcheck(healthcheckUrl)
			state.healthyNodes = healthyNodes
			state.unhealthyNodes = unhealthyNodes

			time.Sleep(healtcheckDelayInMS)
		}
		
	}()

	e.Logger.Fatal(e.Start(port))
}