package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Log visit
	h.db.Exec("INSERT INTO visits (ip_address, user_agent) VALUES (?, ?)",
		r.RemoteAddr, r.UserAgent())

	// Serve the static landing page
	http.ServeFile(w, r, "static/index.html")
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, email, created_at FROM users ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email, createdAt string
		if err := rows.Scan(&id, &name, &email, &createdAt); err != nil {
			continue
		}
		users = append(users, map[string]interface{}{
			"id":         id,
			"name":       name,
			"email":      email,
			"created_at": createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")

	if name == "" || email == "" {
		http.Error(w, "Name and email required", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "created",
	})
}
