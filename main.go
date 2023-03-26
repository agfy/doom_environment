package main

import (
	"github.com/agfy/doom-environment/environment"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	env, err := environment.Create(2)
	if err != nil {
		println(err.Error())
	}
	time.Sleep(10 * time.Second)
	err = env.Start()
	if err != nil {
		println(err.Error())
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	_ = <-sigc
}
