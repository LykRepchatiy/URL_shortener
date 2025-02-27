package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	utils "url_shortener/internal/handler_utils"

	"github.com/jackc/pgx/v5"
)

func Post(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, _ := pgx.Connect(ctx, "postgres://postgres:perlovka14@localhost:5432/patchesj")
	defer conn.Close(ctx)
	sql_create_tabe := "CREATE TABLE IF NOT EXISTS url_data (id SERIAL PRIMARY KEY, short_url TEXT NOT NULL UNIQUE, url TEXT NOT NULL);"
	sql_insert := "INSERT INTO url_data (short_url, url) VALUES ($1, $2)"
	_, _ = conn.Exec(ctx, sql_create_tabe)
	body, _ := io.ReadAll(r.Body)
	var request utils.HTPPModel
	json.Unmarshal(body, &request)
	if request.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	if !utils.IsValidURL(request.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	short_url := utils.ShortURL(request.URL)
	_, err := conn.Exec(ctx, sql_insert, short_url, request.URL)
	if err != nil {
		sql_select_short := "SELECT short_url from url_data WHERE short_url=$1"
		var sql_short_url string
		err = conn.QueryRow(ctx, sql_select_short, short_url).Scan(&sql_short_url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		exist_error := "Data already exist " + sql_short_url
		http.Error(w, errors.New(exist_error).Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(utils.HTPPModel{URL: short_url}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Data saved successfully. Short URL: " + short_url)
}

func Get(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, _ := pgx.Connect(ctx, "postgres://postgres:perlovka14@localhost:5432/patchesj")
	defer conn.Close(ctx)
	short_url := r.URL.Query().Get("short_url")
	sql_select_origin := "SELECT url FROM url_data WHERE short_url=$1"
	var sql_origin_url string
	err := conn.QueryRow(ctx, sql_select_origin, short_url).Scan(&sql_origin_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(utils.HTPPModel{URL: sql_origin_url}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Data fetched successfully. Original URL: " + sql_origin_url)
}
