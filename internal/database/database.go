package database

import (
	"context"
	"errors"
	"url_shortener/internal/service"

	"github.com/jackc/pgx/v5"
)

const (
	Sql_create_table  = "CREATE TABLE IF NOT EXISTS url_data (id SERIAL PRIMARY KEY, short_url TEXT NOT NULL UNIQUE, url TEXT NOT NULL UNIQUE);"
	sql_insert        = "INSERT INTO url_data (short_url, url) VALUES ($1, $2)"
	sql_select_origin = "SELECT url FROM url_data WHERE short_url=$1"
	sql_select_short  = "SELECT short_url FROM url_data WHERE short_url=$1"
)

func checkMatch(DBConn *pgx.Conn, ctx context.Context,
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

func DBPush(DBConn *pgx.Conn, short_url string, request service.HTPPModel) error {
	ctx := context.Background()
	_, err := DBConn.Exec(ctx, sql_insert, short_url, request.URL)
	if err != nil {
		short_url, err := checkMatch(DBConn, ctx, short_url, request.URL)
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

func DBGet(DBConn *pgx.Conn, short_url string) (string, error) {
	ctx := context.Background()
	var sql_origin_url string
	err := DBConn.QueryRow(ctx, sql_select_origin, short_url).Scan(&sql_origin_url)
	if err == pgx.ErrNoRows {
		return "", errors.New("data not found")
	}
	return sql_origin_url, nil
}
