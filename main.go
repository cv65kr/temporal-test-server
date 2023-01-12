package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
)

type App struct {
	resetSignal chan struct{}
}

func NewApp() *App {
	return &App{
		resetSignal: make(chan struct{}),
	}
}

func (a *App) resetControllerHandler(c echo.Context) error {
	a.resetSignal <- struct{}{}
	return c.String(http.StatusOK, "OK")
}

func (a *App) runTestServer() {
	for {
		fmt.Println("[START] Temporal test server")
		cmd := exec.Command("/app/temporal-test-server", "7233", "--enable-time-skipping")
		if err := cmd.Start(); err != nil {
			panic(err)
		}

		<-a.resetSignal
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("failed to kill process")
		}
	}
}

func main() {
	done := make(chan struct{})
	app := NewApp()

	// Run Temporal test server as a child process
	go app.runTestServer()

	// Run HTTP server it will be responsible for handling reset signal
	go func() {
		fmt.Println("[START] HTTP server")
		e := echo.New()
		e.POST("/reset", app.resetControllerHandler)
		e.Logger.Fatal(e.Start(":1323"))
	}()

	<-done
}
