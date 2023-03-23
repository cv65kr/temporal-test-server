package app

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"syscall"

	"github.com/labstack/echo/v4"
)

type App struct {
	resetSignal chan struct{}
	cmd         *exec.Cmd
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

func (a *App) RunTestServer(path *string) {

	for {
		fmt.Println("[START] Temporal test server from " + *path)
		a.cmd = exec.Command(*path, "7233", "--enable-time-skipping")

		stderr, err := a.cmd.StderrPipe()
		if err != nil {
			log.Fatalf("could not get stderr pipe: %v", err)
		}

		stdout, err := a.cmd.StdoutPipe()
		if err != nil {
			log.Fatalf("could not get stdout pipe: %v", err)
		}

		go func() {
			merged := io.MultiReader(stderr, stdout)
			scanner := bufio.NewScanner(merged)
			for scanner.Scan() {
				msg := scanner.Text()
				fmt.Printf("msg: %s\n", msg)
			}
		}()

		if err := a.cmd.Start(); err != nil {
			log.Fatal(err)
		}

		<-a.resetSignal

		if err := a.cmd.Process.Signal(syscall.SIGKILL); err != nil {
			log.Fatal("failed to kill process")
		}
	}
}

func (a *App) Stop() {
	if err := a.cmd.Process.Signal(syscall.SIGKILL); err != nil {
		log.Fatal("failed to kill process")
	}
}
