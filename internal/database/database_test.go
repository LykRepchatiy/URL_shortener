package database

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"url_shortener/internal/service"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"

	"github.com/stretchr/testify/assert"
)

func TestCheckMatch(t *testing.T) {
	ctx := context.Background()

	mockDB, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	t.Run("Match Found", func(t *testing.T) {
		mockDB.ExpectQuery(sql_select_origin).
			WithArgs("8YtucraCrC").
			WillReturnRows(pgxmock.NewRows([]string{"url"}).AddRow("https://example.com"))

		short, err := checkMatch(mockDB, ctx, "8YtucraCrC", "https://example.com")
		assert.NoError(t, err)
		assert.Equal(t, "8YtucraCrC", short)
	})
	mockDB.Close(ctx)

	mockDB, err = pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	t.Run("Data Not Found", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(sql_select_origin)).
			WithArgs("8YtucraCrC").
			WillReturnError(pgx.ErrNoRows)

		short, err := checkMatch(mockDB, ctx, "8YtucraCrC", "https://example.com")
		assert.Error(t, err)
		assert.Equal(t, "data not found", err.Error())
		assert.Empty(t, short)
	})
	mockDB.Close(ctx)

	mockDB, err = pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	t.Run("URL Mismatch", func(t *testing.T) {
		mockDB.ExpectQuery(sql_select_origin).
			WithArgs("8YtucraCrC").
			WillReturnRows(pgxmock.NewRows([]string{"url"}).AddRow("https://different.com"))

		short, err := checkMatch(mockDB, ctx, "8YtucraCrC", "https://example.com/")
		assert.NoError(t, err)
		assert.NotEqual(t, "8YtucraCrC", short)
	})

	mockDB.Close(ctx)
}

func TestDBPush_SuccessfulInsertion(t *testing.T) {
	mockDB, err := pgxmock.NewConn(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Не удалось создать mock-базу данных: %v", err)
	}
	defer mockDB.Close(context.Background())

	mockDB.ExpectExec(sql_insert).
		WithArgs("8YtucraCrC", "https://example.com").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = DBPush(mockDB, "8YtucraCrC", service.HTPPModel{URL: "https://example.com"})
	assert.NoError(t, err)

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены все ожидания: %s", err)
	}
}

func TestDBPush_CheckMatchError(t *testing.T) {
	mockDB, err := pgxmock.NewConn(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Не удалось создать mock-базу данных: %v", err)
	}
	defer mockDB.Close(context.Background())

	mockDB.ExpectExec(sql_insert).
		WithArgs("8YtucraCrC", "https://example.com").
		WillReturnError(errors.New("insert error"))

	mockDB.ExpectQuery(sql_select_origin).
		WithArgs("8YtucraCrC").
		WillReturnError(pgx.ErrNoRows)

	err = DBPush(mockDB, "8YtucraCrC", service.HTPPModel{URL: "https://example.com"})
	assert.Error(t, err)
	assert.Equal(t, "data not found", err.Error())

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены все ожидания: %s", err)
	}
}

func TestDBGet_Success(t *testing.T) {
	ctx := context.Background()
	mockDB, err := pgxmock.NewConn(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Не удалось создать mock-базу данных: %v", err)
	}
	defer mockDB.Close(ctx)

	mockDB.ExpectQuery(sql_select_origin).
		WithArgs("8YtucraCrC").
		WillReturnRows(pgxmock.NewRows([]string{"url"}).AddRow("https://example.com"))

	originURL, err := DBGet(mockDB, "8YtucraCrC")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originURL)

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены все ожидания: %s", err)
	}
}

func TestDBGet_DataNotFound(t *testing.T) {
	ctx := context.Background()
	mockDB, err := pgxmock.NewConn(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Не удалось создать mock-базу данных: %v", err)
	}
	defer mockDB.Close(ctx)

	mockDB.ExpectQuery(sql_select_origin).
		WithArgs("8YtucraCrC").
		WillReturnError(pgx.ErrNoRows)

	originURL, err := DBGet(mockDB, "8YtucraCrC")
	assert.Error(t, err)
	assert.Equal(t, "", originURL)
	assert.Equal(t, "data not found", err.Error())

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены все ожидания: %s", err)
	}
}