package main

// Handles signals
import (
	"os"
	"os/signal"
	"syscall"
)

var signalChan = make(chan os.Signal, 1)

func handleSignals() {
	signal.Notify(signalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	go func() {
		sig := <-signalChan
		switch sig {
		case syscall.SIGINT:
			// handle SIGINT
			println("SIGINT recieved!")
			os.Exit(0)
		case syscall.SIGTERM:
			println("SIGTERM recieved!")
			os.Exit(0)			
		}

	}()
}
