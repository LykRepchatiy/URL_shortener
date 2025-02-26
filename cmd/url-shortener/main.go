package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Request struct {
	URL string `json:"url"`
}

//TO DO: flag

// func init() {
// 	flag.BoolVar(p, "")
// }

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

func shortURL(input string) string {
	hash := sha256.Sum256([]byte(input))
	var builder strings.Builder
	builder.Grow(10)
	for i := 0; i < 10; i++ {
		index := int(hash[i]) % len(alphabet)
		builder.WriteByte(alphabet[index])
	}
	return builder.String()
}

func writeToDatabase(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5433/postgres")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(ctx)
	sql_create_tabe := `CREATE TABLE IF NOT EXISTS url_data (id SERIAL PRIMARY KEY, short_url TEXT NOT NULL UNIQUE, url TEXT NOT NULL);`
	sql_insert := "INSERT INTO url_data (short_url, url) VALUES ($1, $2)"
	_, err = conn.Exec(ctx, sql_create_tabe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var request Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	short_url := shortURL(request.URL)
	_, err = conn.Exec(ctx, sql_insert, short_url, request.URL)
	if err != nil {
		sql_select_short := "SELECT short_url from url_data WHERE short_url = '($1)'"
		row, err := conn.Query(ctx, sql_select_short, short_url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(row)
		defer row.Close()
		var sql_short_url string
		for row.Next() {
			err := row.Scan(&sql_short_url)
			if err != nil {
				log.Println("4")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		exist_error := "Data already exist " + sql_short_url
		log.Println(sql_short_url)
		http.Error(w, errors.New(exist_error).Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	r := chi.NewRouter()
	r.Post("/database", writeToDatabase)
	http.ListenAndServe(":8080", r)
}
