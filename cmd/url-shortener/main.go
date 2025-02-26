package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Request struct {
	URL string `json:"url"`
}

func writeToDatabase(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:perlovka14@localhost:5432/patchesj")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(ctx)
	sql_create_tabe := `CREATE TABLE IF NOT EXISTS url_data ( id SERIAL PRIMARY KEY, url TEXT NOT NULL);`
	sql_insert := "INSERT INTO url_data (url) VALUES ($1)"
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
	_, err = conn.Exec(ctx, sql_insert, request.URL)
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
