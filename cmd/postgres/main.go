package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var (
	defaultConnStr = "postgresql://postgres:password@localhost/postgres?sslmode=disable"
)

func main() {

	connStr, ok := os.LookupEnv("CONNECTION_STRING")
	if !ok {
		log.Println("CONNECTION_STRING not set. Using default.")
		connStr = defaultConnStr
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, "USE POST TO UPLOAD DATA")
		} else if r.Method == "POST" {
			var data Word
			body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
			if err != nil {
				panic(err)
			}
			if err := r.Body.Close(); err != nil {
				panic(err)
			}
			w.Header().Set("Content-Type", "application/json;   charset=UTF-8")
			if err := json.Unmarshal(body, &data); err != nil {
				panic(err)
			}
			query := fmt.Sprintf("INSERT INTO words VALUES('%s',%d);", data.Name, data.Count)
			_, err = db.Query(query)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "{\"status\":\"success\"}")
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
