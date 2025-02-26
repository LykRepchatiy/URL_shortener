package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"io"
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
	conn, err := pgx.Connect(ctx, "postgres://postgres:perlovka14@localhost:5432/patchesj")
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
	_, err = conn.Exec(ctx, sql_insert, shortURL(request.URL), request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	r := chi.NewRouter()
	r.Post("/database", writeToDatabase)
	http.ListenAndServe(":8080", r)
}
