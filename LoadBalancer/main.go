package main

import (
	"fmt"
	"flag"

	"github.com/labstack/echo/v4"
	
	"LoadBalancer/handlers"
)

func main() {
	portNumber := flag.String("port", "8080", "LoadBalancer Port")
	flag.Parse()
	port := fmt.Sprintf(":%s", *portNumber)

	e := echo.New()

	e.GET("/register", handlers.GetRegisterHandler)
	e.POST("/*", handlers.PostEchoHandler)

	go handlers.InitiateHealthcheck()

	e.Logger.Fatal(e.Start(port))
}