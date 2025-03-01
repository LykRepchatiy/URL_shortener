package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock"
)

func TestPostDB(t *testing.T) {
	// TO OD последний тест не проходит, доделать
	conn, _ := pgxmock.NewConn()

	r := &Router{
		PG: conn,
	}

	t.Run("Отсутствует URL", func(t *testing.T) {
		reqBody := []byte(`{"URL":""}`)
		req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		r.PostDB(w, req)
		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("ожидался статус %d, получен %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("Недопустимый URL", func(t *testing.T) {
		reqBody := []byte(`{"URL":"invalid"}`)
		req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		r.PostDB(w, req)
		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("ожидался статус %d, получен %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("Ошибка DBPush", func(t *testing.T) {
		reqBody := []byte(`{"URL":"http://example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		r.PostDB(w, req)
		resp := w.Result()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("ожидался статус %d, получен %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})


	t.Run("Успешное выполнение", func(t *testing.T) {
		reqBody := []byte(`{"URL":"http://example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		
		r.PostDB(w, req)
		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("ожидался статус %d, получен %d", http.StatusOK, resp.StatusCode)
		}
	})
}
