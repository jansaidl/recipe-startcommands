package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	ctx := ContextWithSigterm(context.Background())

	var port int
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.Parse()

	fmt.Printf("Listening on port %d\n", port)
	var count uint64
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World! port: %d - %d", port, count)
	})}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				fmt.Printf("Error starting server: %s\n", err)
				return
			}
		}
		fmt.Printf("Server stopped on port: %d\n", port)
	}()

	for {
		atomic.AddUint64(&count, 1)
		select {
		case <-time.After(time.Minute):
			fmt.Printf("Hello World! port: %d - %d\n", port, count)
			continue
		case <-ctx.Done():
			fmt.Printf("Stopping server on port: %d\n", port)
			server.Close()
		}
		break
	}
	wg.Wait()
}

func ContextWithSigterm(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		<-interrupt
		cancel()
	}()
	return ctx
}
