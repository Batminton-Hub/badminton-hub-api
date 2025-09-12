package util

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var srv *http.Server

func ShutdownServer() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Forced shutdown:", err)
	}

	fmt.Println("Server exited gracefully")
}

func HttpServer(handler http.Handler) *http.Server {
	config := LoadConfig()
	srv = &http.Server{
		Addr:    config.ServerPort,
		Handler: handler,
	}

	fmt.Println("Server port:", config.ServerPort)
	return srv
}

func RunServer(handler http.Handler) {
	srv := HttpServer(handler)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Listen error:", err)
		}
	}()
}
