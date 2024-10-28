package main

import (
	"fmt"
	"gameApp/config"
	"gameApp/scheduler"
	"os"
	"os/signal"
	"time"
)

func main() {
	// TODO - read config path from command line

	cfg := config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n", cfg)

	done := make(chan bool)

	go func() {
		sch := scheduler.New()
		sch.Start(done)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("\n recieved interrupt signal, shutting down gracefully...")
	done <- true
	time.Sleep(cfg.Application.GraceFullShutdownTimeout)

}
