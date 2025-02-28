package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"url_shortener/internal/cache"
	"url_shortener/internal/database"
	"url_shortener/internal/service"
)

func (rout *Router) PostDB(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var request service.HTPPModel
	json.Unmarshal(body, &request)
	if request.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	if !service.IsValidURL(request.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	short_url := service.ShortURL(request.URL)
	err := database.DBPush(rout.PG, short_url, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(service.HTPPModel{URL: short_url}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Data saved successfully. Short URL: " + short_url)
}

func (rout *Router) GetDB(w http.ResponseWriter, r *http.Request) {
	short_url := r.URL.Query().Get("short_url")
	sql_origin_url, err := database.DBGet(rout.PG, short_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(service.HTPPModel{URL: sql_origin_url}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Data fetched successfully. Original URL: " + sql_origin_url)
}

func (rout *Router) PostCache(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var request service.HTPPModel
	json.Unmarshal(body, &request)
	if request.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	if !service.IsValidURL(request.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	short_url := service.ShortURL(request.URL)
	err := cache.PushCache(short_url, request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Data saved successfully. Short URL: " + short_url)
}

func (rout *Router) GetCache(w http.ResponseWriter, r *http.Request) {
	short_url := r.URL.Query().Get("short_url")
	origin_url, err := cache.GetCache(short_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(service.HTPPModel{URL: origin_url}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Data fetched successfully. Original URL: " + origin_url)
}
