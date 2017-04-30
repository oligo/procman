package main

// Handles signals
import (
	"os"
	"os/signal"
	"syscall"
)

type SigHandler func() error

var signalChan = make(chan os.Signal, 1)
var signalHandlers = make(map[os.Signal]SigHandler)

func registerHandler(sig os.Signal, handler SigHandler) {
	signalHandlers[sig] = handler
}

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
			if handler, hasHandler := signalHandlers[sig]; hasHandler {
				if handler() != nil {
					os.Exit(1)
				}
			}
			os.Exit(0)
		case syscall.SIGTERM:
			if handler, hasHandler := signalHandlers[sig]; hasHandler {
				if handler() != nil {
					os.Exit(1)
				}
			}
			os.Exit(0)
		}
	}()
}
