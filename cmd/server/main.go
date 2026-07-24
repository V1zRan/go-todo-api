// main.go
//
// Точка входа приложения.
// Здесь настраиваются маршруты и запускается HTTP-сервер.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-api/internal/config"
	"todo-api/internal/handlers"
	"todo-api/internal/storage/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	postgresStorage, err := postgres.NewStorage(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Home page")
		fmt.Fprintf(w, "Method: %s\n", r.Method)
		fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
		fmt.Fprintf(w, "Host: %s\n", r.Host)
	})

	http.HandleFunc("/tasks", handlers.TasksHandler(postgresStorage))
	http.HandleFunc("/tasks/", handlers.HandleTaskByID(postgresStorage))

	fmt.Println("Server is running on http://localhost:8080")
	fmt.Println("Page http://localhost:8080/tasks")

	// Запуск сервера
	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	//Ждёт, пока signalCtx не сообщит, что он завершён
	<-signalCtx.Done()

	fmt.Println("Получен сигнал завершения")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Ошибка завершения HTTP-сервера:", err)
	}

	fmt.Println("HTTP-сервер завершён")

	if err := postgresStorage.Close(); err != nil {
		fmt.Println("Ошибка закрытия PostgreSQL:", err)
	} else {
		fmt.Println("Соединение с PostgreSQL закрыто")
	}

}
