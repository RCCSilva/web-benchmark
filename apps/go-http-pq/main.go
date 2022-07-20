package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const jsonContentType = "application/json"

type response struct {
	Status string
}

type User struct {
	Id   int
	Name string
	Age  int
}

type BenchmarkServer struct {
	db *sql.DB
	http.Handler
}

func main() {
	godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	db.SetMaxOpenConns(100)

	if err != nil {
		log.Fatal(err)
	}

	httpServer := NewBenchmarkServer(db)
	log.Fatal(http.ListenAndServe(":3000", httpServer))
}

func (p *BenchmarkServer) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	// Create
	var userId int
	err := p.db.QueryRow("INSERT INTO users (name, age) VALUES ('Test!', 1) RETURNING id;").Scan(&userId)

	if err != nil {
		WriteError(w, err)
		return
	}

	// Read
	var userDb User
	err = p.db.QueryRow(
		"SELECT id, name, age FROM users WHERE id = $1;", userId,
	).Scan(&userDb.Id, &userDb.Name, &userDb.Age)
	if err != nil {
		WriteError(w, err)
		return
	}

	// Update
	userDb.Age += 1
	err = p.db.QueryRow(
		"UPDATE users SET name = $2, age = $3 WHERE id = $1 RETURNING id, name, age;",
		userDb.Id, userDb.Name, userDb.Age,
	).Scan(&userDb.Id, &userDb.Name, &userDb.Age)

	if err != nil {
		WriteError(w, err)
		return
	}

	// Delete
	_, err = p.db.Exec("DELETE FROM users WHERE id = $1", userDb.Id)
	if err != nil {
		WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(userDb)
}

func NewBenchmarkServer(db *sql.DB) *BenchmarkServer {
	server := new(BenchmarkServer)
	server.db = db

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(server.index))

	server.Handler = router

	return server
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
