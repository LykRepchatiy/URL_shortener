package database

import (
	"context"
	"errors"
	"url_shortener/internal/service"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type QueryRower interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Close(ctx context.Context) error
}

type DataBase struct {}

type DB interface {
	СheckMatch(DBConn QueryRower, ctx context.Context,
		short_url, URL string) (string, error)
	DBPush(DBConn QueryRower, short_url string, request service.HTTPModel) error
	DBGet(DBConn QueryRower, short_url string) (string, error)
}

const (
	Sql_create_table  = "CREATE TABLE IF NOT EXISTS url_data (id SERIAL PRIMARY KEY, short_url TEXT NOT NULL UNIQUE, url TEXT NOT NULL UNIQUE);"
	Sql_insert        = "INSERT INTO url_data (short_url, url) VALUES ($1, $2)"
	sql_select_origin = "SELECT url FROM url_data WHERE short_url=$1"
	sql_select_short   = "SELECT short_url FROM url_data WHERE short_url=$1"
)

func (db DataBase) СheckMatch(DBConn QueryRower, ctx context.Context,
	short_url, URL string) (string, error) {
	var sql_origin_url string
	err := DBConn.QueryRow(ctx, sql_select_origin, short_url).Scan(&sql_origin_url)
	if err == pgx.ErrNoRows {
		return "", errors.New("data not found")
	}
	if sql_origin_url != URL {
		tmp := URL[:len(URL)-1]
		newShort := service.ShortURL(tmp)
		short_url = newShort
	}
	return short_url, nil
}

func (db DataBase) DBPush(DBConn QueryRower, short_url string, request service.HTTPModel) error {
	ctx := context.Background()
	_, err := DBConn.Exec(ctx, Sql_insert, short_url, request.URL)
	if err != nil {
		short_url, err := db.СheckMatch(DBConn, ctx, short_url, request.URL)
		if err != nil {
			return err
		}
		var sql_short_url string
		err = DBConn.QueryRow(ctx, sql_select_short, short_url).Scan(&sql_short_url)
		if err != nil {
			return err
		}
		exist_error := "Data already exist " + sql_short_url
		return errors.New(exist_error)
	}
	return nil
}

func (db DataBase) DBGet(DBConn QueryRower, short_url string) (string, error) {
	ctx := context.Background()
	var sql_origin_url string
	err := DBConn.QueryRow(ctx, sql_select_origin, short_url).Scan(&sql_origin_url)
	if err == pgx.ErrNoRows {
		return "", errors.New("data not found")
	}
	return sql_origin_url, nil
}
