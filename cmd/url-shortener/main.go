package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"url_shortener/internal/cache"
	"url_shortener/internal/database"
	"url_shortener/internal/handlers"

	"github.com/jackc/pgx/v5"
)

var (
	postgre bool
)

func init() {
	flag.BoolVar(&postgre, "p", false, "use postgre for store URLs")
	flag.Parse()
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	signalChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	r := handlers.NewRouter()
	if postgre {
		ctx := context.Background()
		DBConn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
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
	} else {
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
