package main

import (
	"flag"
	"fmt"
	"github/cv65kr/temporal-test-server/internal/app"

	"github.com/labstack/echo/v4"
)

func main() {
	path := flag.String("path", "/app/temporal-test-server", "Path to temporal test server")
	flag.Parse()

	done := make(chan struct{})
	app := app.NewApp()

	// Run Temporal test server as a child process
	go app.RunTestServer(*path)

	// Run HTTP server it will be responsible for handling reset signal
	go func() {
		fmt.Println("[START] HTTP server")
		e := echo.New()
		e.POST("/reset", app.ResetControllerHandler)
		e.Logger.Fatal(e.Start(":1323"))
	}()

	<-done
}
