package main

import (
	"database/sql"
	"fmt"
	"go-blog/handlers"
	"go-blog/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Database configuration from env vars with defaults
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "bloguser")
	dbPass := getEnv("DB_PASSWORD", "blogpass")
	dbName := getEnv("DB_NAME", "blogdb")
	serverPort := getEnv("PORT", "8080")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	// Connect to database with retry
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Waiting for database... attempt %d/10: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup store and handlers
	store := models.NewPostStore(db)
	userStore := models.NewUserStore(db)
	sessionStore := models.NewSessionStore(db)
	
	blogHandler := handlers.NewBlogHandler(store)
	authHandler := handlers.NewAuthHandler(userStore, sessionStore)

	// Setup router
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// Public routes
	r.HandleFunc("/", blogHandler.HomePage).Methods("GET")
	r.HandleFunc("/post/{slug}", blogHandler.PostPage).Methods("GET")
	r.HandleFunc("/login", authHandler.LoginPage).Methods("GET")
	r.HandleFunc("/login", authHandler.LoginPost).Methods("POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("GET", "POST")

	// Admin routes
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authHandler.AuthMiddleware)
	adminRouter.HandleFunc("", blogHandler.AdminPage).Methods("GET")
	adminRouter.HandleFunc("/new", blogHandler.AdminNewPage).Methods("GET")
	adminRouter.HandleFunc("/edit/{id:[0-9]+}", blogHandler.AdminEditPage).Methods("GET")
	adminRouter.HandleFunc("/create", blogHandler.CreatePost).Methods("POST")
	adminRouter.HandleFunc("/update/{id:[0-9]+}", blogHandler.UpdatePost).Methods("POST")
	adminRouter.HandleFunc("/delete/{id:[0-9]+}", blogHandler.DeletePost).Methods("POST")

	addr := ":" + serverPort
	log.Printf("Server running at http://localhost%s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		"migrations/001_create_posts.sql",
		"migrations/002_create_users.sql",
	}

	for _, file := range migrations {
		migration, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(migration)); err != nil {
			return fmt.Errorf("execute migration %s: %w", file, err)
		}
	}

	log.Println("Migrations executed successfully")
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
