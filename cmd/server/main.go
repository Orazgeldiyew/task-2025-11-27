// Link Status Checker API.
//
// Сервис принимает списки ссылок, асинхронно проверяет их доступность
// и позволяет сформировать PDF-отчёт по ранее созданным наборам ссылок.
//
// Тестовое задание: веб-сервис без внешней инфраструктуры (DB, Docker и т.д.).
//
// @title           Link Status Checker API
// @version         1.0
// @description     Сервис для проверки доступности интернет-ссылок и генерации PDF-отчётов по наборам ссылок.
// @host            localhost:8080
// @BasePath        /
// @schemes         http
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task/internal/service"
	"task/internal/storage"
	"task/internal/transport"
)

func main() {

	st := storage.NewFileStorage("data/state.json")

	mgr, err := service.NewManager(st)
	if err != nil {
		log.Fatalf("cannot init manager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mgr.StartWorkers(ctx, 4)

	router := transport.NewRouter(mgr)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Server is running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown error: %v", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}
}
