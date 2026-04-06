package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"boatly/internal/database"
	"boatly/internal/handlers"
)

func main() {
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	// Initialize database
	dbPath := filepath.Join(wd, "data", "boatly.db")
	db, err := database.Init(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Setup handlers
	h := handlers.New(db)

	// Routes
	http.HandleFunc("/", h.Index)
	http.HandleFunc("/api/health", h.Health)
	http.HandleFunc("/api/users", h.Users)
	http.HandleFunc("/api/users/create", h.CreateUser)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(wd, "static")))))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3050"
	}

	log.Printf("🚀 Boatly server starting on port %s", port)
	log.Printf("📁 Database: %s", dbPath)
	log.Printf("🌐 http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
