package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

var db *sql.DB

func response(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)

	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func handler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		users := []User{
			{"Gabriel", 26},
			{"Marianinha", 24},
			{"Nana", 8},
		}

		response(res, http.StatusOK, users)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := map[string]any{
			"status": "up",
			"time":   time.Now().Unix(),
		}

		response(w, http.StatusOK, data)
	}
}

func main() {
	var port string

	if port = os.Getenv("PORT"); port == "" {
		log.Fatal(errors.New("PORT is not defined"))
	}

	var err error

	db, err = sql.Open("postgres", os.Getenv("URL_CONNECT"))
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	log.Println("PING: ", db.Ping())

	mux := http.NewServeMux()

	mux.HandleFunc("/users", handler)
	mux.HandleFunc("/health", handleHealth)

	log.Printf("server stated on http://localhost:%s\n", port)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalln(err)
	}
}
