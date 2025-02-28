package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"url_shortener/internal/database"
	"url_shortener/internal/handlers"

	"github.com/jackc/pgx/v5"
)

var (
	postgre, cache bool
)

func init() {
	flag.BoolVar(&postgre, "postgre", false, "use postgre for store URLs")
	flag.BoolVar(&cache, "cache", false, "use cache for store URLs")
	flag.Parse()
	if !postgre && !cache {
		postgre = true
	}
	if cache && postgre {
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		logger.Error("Both postgre and cache flags are set. Please choose only one.")
		os.Exit(1)
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	signalChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	r := handlers.NewRouter()
	if postgre {
		ctx := context.Background()
		DBConn, err := pgx.Connect(ctx, "postgres://postgres:perlovka14@localhost:5432/patchesj")
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
	} else if cache {

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
