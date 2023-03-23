package main

import (
	"flag"
	"fmt"
	"github/cv65kr/temporal-test-server/internal/app"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

func main() {
	path := flag.String("path", "/app/temporal-test-server", "Path to temporal test server")
	flag.Parse()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	app := app.NewApp()

	// Run Temporal test server as a child process
	go func() {
		app.RunTestServer(path)
	}()

	// Run HTTP server it will be responsible for handling reset signal
	go func() {
		fmt.Println("[START] HTTP server")
		e := echo.New()
		e.POST("/reset", app.ResetControllerHandler)
		e.Logger.Fatal(e.Start(":1323"))
	}()

	<-stopCh
	app.Stop()
}
