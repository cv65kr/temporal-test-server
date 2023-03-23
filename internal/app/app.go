package app

import (
	"bufio"
	"fmt"
	"log"
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

func (a *App) ResetControllerHandler(c echo.Context) error {
	a.resetSignal <- struct{}{}
	return c.String(http.StatusOK, "OK")
}

func (a *App) RunTestServer(path string) {
	for {
		fmt.Println("[START] Temporal test server")
		cmd := exec.Command(path, "7233", "--enable-time-skipping")
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(stderr)
		line, err := reader.ReadString('\n')
		for err == nil {
			fmt.Println(line)
			line, err = reader.ReadString('\n')
		}

		<-a.resetSignal
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process")
		}
	}
}
