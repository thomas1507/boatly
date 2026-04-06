package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
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

	// Count visits
	var visitCount int
	h.db.QueryRow("SELECT COUNT(*) FROM visits").Scan(&visitCount)

	// Count users
	var userCount int
	h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)

	data := struct {
		VisitCount int
		UserCount  int
		ServerTime string
	}{
		VisitCount: visitCount,
		UserCount:  userCount,
		ServerTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Boatly - Welcome Aboard!</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            color: white;
        }
        .container {
            text-align: center;
            padding: 2rem;
            max-width: 600px;
        }
        h1 {
            font-size: 4rem;
            margin-bottom: 1rem;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        .emoji {
            font-size: 5rem;
            margin-bottom: 1rem;
        }
        .stats {
            background: rgba(255,255,255,0.2);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 2rem;
            margin: 2rem 0;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 2rem;
        }
        .stat h3 {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
        }
        .stat p {
            opacity: 0.9;
        }
        .server-time {
            margin-top: 2rem;
            opacity: 0.8;
            font-size: 0.9rem;
        }
        .tech-stack {
            margin-top: 2rem;
            display: flex;
            justify-content: center;
            gap: 1rem;
            flex-wrap: wrap;
        }
        .badge {
            background: rgba(255,255,255,0.3);
            padding: 0.5rem 1rem;
            border-radius: 20px;
            font-size: 0.85rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="emoji">🚢</div>
        <h1>Boatly</h1>
        <p>Welcome aboard! Your Go + SQLite app is running!</p>
        
        <div class="stats">
            <div class="stats-grid">
                <div class="stat">
                    <h3>{{.VisitCount}}</h3>
                    <p>Total Visits</p>
                </div>
                <div class="stat">
                    <h3>{{.UserCount}}</h3>
                    <p>Registered Users</p>
                </div>
            </div>
        </div>
        
        <div class="tech-stack">
            <span class="badge">Go</span>
            <span class="badge">SQLite</span>
            <span class="badge">Caddy</span>
        </div>
        
        <div class="server-time">
            Server time: {{.ServerTime}}
        </div>
    </div>
</body>
</html>
	`

	t := template.Must(template.New("index").Parse(tmpl))
	t.Execute(w, data)
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
