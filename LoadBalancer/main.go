package main

import (
	"fmt"
	"flag"
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

var state loadBalancerStatus

func main() {
	portNumber := flag.String("port", "8080", "LoadBalancer Port")
	flag.Parse()
	port := fmt.Sprintf(":%s", *portNumber)

	e := echo.New()

	e.GET("/register", getRegisterHandler)
	e.POST("/echo", postEchoHandler)

	e.Logger.Fatal(e.Start(port))
}