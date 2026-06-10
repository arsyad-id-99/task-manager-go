package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/arsyad-id-99/task-manager-go/internal/handler"
	"github.com/arsyad-id-99/task-manager-go/internal/middleware"
	"github.com/arsyad-id-99/task-manager-go/internal/repository"
)

func main() {
	_ = godotenv.Load()

	// Koneksi ke PostgreSQL
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Inisialisasi layer
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	authHandler := handler.NewAuthHandler(userRepo)
	taskHandler := handler.NewTaskHandler(taskRepo)

	// Setup router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Public routes
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	// Protected routes — JWT
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/tasks", taskHandler.List)
		r.Post("/tasks", taskHandler.Create)
		r.Get("/tasks/{id}", taskHandler.Detail)
		r.Patch("/tasks/{id}/status", taskHandler.UpdateStatus)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
