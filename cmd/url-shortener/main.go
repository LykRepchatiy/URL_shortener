package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"url_shortener/internal/cache"
	"url_shortener/internal/database"
	"url_shortener/internal/handlers"

	"github.com/jackc/pgx/v4"
)

var (
	postgre bool
)

type Env struct {
	databaseUrl string
	storage     string
}

func initEnv() Env {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	storage := os.Getenv("STORAGE")
	if storage == "" {
		log.Fatal("STORAGE is not set")
	}
	return Env{
		databaseUrl: databaseUrl,
		storage:     storage}
}

func main() {
	env := initEnv()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	signalChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	c := cache.Init()
	db := database.DataBase{}
	r := handlers.NewRouter(db, c)
	if env.storage == "postgre" {
		ctx := context.Background()
		DBConn, err := pgx.Connect(ctx, env.databaseUrl)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		_, err = DBConn.Exec(ctx, database.Sql_create_table)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		go func() {
			errChan <- r.StartDB(DBConn)
		}()
	} else if env.storage == "cache" {
		r.Cache = cache.Init()
		go func() {
			errChan <- r.StartCache()
		}()
	}

	select {
	case <-signalChan:
		r.Finish()
	case <-errChan:
		r.Finish()
	}
	close(signalChan)
	close(errChan)
}
