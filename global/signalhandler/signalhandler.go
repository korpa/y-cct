package signalhandler

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Init(ctx context.Context) (chan os.Signal, context.CancelFunc) {

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// Starting signal handler
	go signalHandler(cancel, c)

	return c, cancel
}

func signalHandler(cancel context.CancelFunc, c chan os.Signal) {
	for {
		s := <-c
		switch s {
		// kill -SIGHUP XXXX
		case syscall.SIGHUP:
			slog.Warn("hungup")
			cancel()
			// exit_chan <- true

		// kill -SIGINT XXXX or Ctrl+c
		case syscall.SIGINT:
			slog.Warn("Ctrl-C received")
			cancel()
			// exit_chan <- true

		// kill -SIGTERM XXXX
		case syscall.SIGTERM:
			slog.Warn("force stop kill 9")
			cancel()
			// exit_chan <- true

		// kill -SIGQUIT XXXX
		case syscall.SIGQUIT:
			slog.Warn("stop and core dump")
			cancel()
			// exit_chan <- true

		default:
			slog.Warn("Unknown signal.")
			cancel()
			// exit_chan <- true
		}
	}
}
