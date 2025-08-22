package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	todos   []Todo
	nextID  = 1
	todoMux sync.Mutex
)

func main() {
	// this is api routes but it not working
	http.HandleFunc("/api/todos", handleTodos)
	http.HandleFunc("/api/todos/", handleTodo)
	//  this serve the frontend file
	http.HandleFunc("/", serveFrontend)

	fmt.Println("🍀 Lucky's To-Do running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "c:/Users/lucky/OneDrive/Desktop/lucky-todo/Untitled-2.html") // must be in the same folder where you run the app
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}
	switch r.Method {
	case "GET":
		getTodos(w)
	case "POST":
		createTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	idStr := r.URL.Path[len("/api/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		updateTodo(w, r, id)
	case "DELETE":
		deleteTodo(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTodos(w http.ResponseWriter) {
	todoMux.Lock()
	defer todoMux.Unlock()
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Text == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	todoMux.Lock()
	todo := Todo{ID: nextID, Text: req.Text, Completed: false, CreatedAt: time.Now()}
	todos = append(todos, todo)
	nextID++
	todoMux.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		Text      *string `json:"text,omitempty"`
		Completed *bool   `json:"completed,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	todoMux.Lock()
	defer todoMux.Unlock()
	for i, t := range todos {
		if t.ID == id {
			if req.Text != nil {
				todos[i].Text = *req.Text
			}
			if req.Completed != nil {
				todos[i].Completed = *req.Completed
			}
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func deleteTodo(w http.ResponseWriter, id int) {
	todoMux.Lock()
	defer todoMux.Unlock()
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
