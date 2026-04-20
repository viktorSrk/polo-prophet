package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}

func (s *Server) Start(addr string) {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server running on http://localhost" + addr)
	http.ListenAndServe(addr, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome!"})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
}
