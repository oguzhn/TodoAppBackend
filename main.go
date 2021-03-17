package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var TodoSvc *MockTodoService

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	TodoSvc = NewMockTodoService()
	mux := http.NewServeMux()

	mux.Handle("/todos", commonHandlers(todoHandler))
	mux.Handle("/todos/", commonHandlers(todoHandler))

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	key := ""
	if len(parts) > 2 {
		key = parts[2]
	}

	switch r.Method {
	case "GET":
		todos, err := TodoSvc.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(todos)

	case "POST":
		if len(key) > 0 {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		todo := Todo{}
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			http.Error(w, err.Error(), 422)
			return
		}
		err = TodoSvc.Save(&todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
