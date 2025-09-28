package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var srv *http.Server

func ShutdownServer(closeFunc ...func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Wait for shutdown signal
	<-quit
	fmt.Println("Shutdown signal received")

	ctx, cancel := InitConText(30 * time.Second)
	defer cancel()

	// Close services
	for _, close := range closeFunc {
		close()
	}

	// Shutdown server
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

	fmt.Println("Server port", config.ServerPort)
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
