package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	ID int `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Done bool `json:"done"`
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todos []Todo

	db.Find(&todos)

	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var todo Todo

	db.First(&todo, vars["id"])

	json.NewEncoder(w).Encode(todo)
}

func toggleTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var todo Todo

	db.First(&todo, vars["id"])
	todo.Done = !todo.Done
	db.Save(&todo)

	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var todo Todo
	var todoUpdates Todo

	json.NewDecoder(r.Body).Decode(&todoUpdates)

	db.First(&todo, vars["id"])
	db.Model(&todo).Updates(todoUpdates)

	json.NewEncoder(w).Encode(todo)

}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var todo Todo

	db.First(&todo, vars["id"])
	db.Delete(&todo)

	json.NewEncoder(w).Encode(Todo{})
}


func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)

	db.Create(&todo)
	json.NewEncoder(w).Encode(todo)
}

var db *gorm.DB
var err error

func main() {
	godotenv.Load()

	dsn := os.Getenv("DSN")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Todo{})


	r := mux.NewRouter()

	r.HandleFunc("/api/todos", getTodos).Methods("GET")
	r.HandleFunc("/api/todos", createTodo).Methods("POST")
	r.HandleFunc("/api/todos/{id:[0-9]+}", getTodo).Methods("GET")
	r.HandleFunc("/api/todos/toggle/{id}", toggleTodo).Methods("POST")
	r.HandleFunc("/api/todos/{id:[0-9]+}", updateTodo).Methods("PUT")
	r.HandleFunc("/api/todos/{id:[0-9]+}", deleteTodo).Methods("DELETE")

	log.Fatal(http.ListenAndServe("localhost:" + os.Getenv("PORT"), r))
}