package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"url_shortener/internal/mocks"
	"url_shortener/internal/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedDbcon := mocks.NewMockDB(ctrl)

	r := &Router{dbCon: mockedDbcon}

	t.Run("Success", func(t *testing.T) {
		requestBody := service.HTTPModel{URL: "https://example.com"}
		jsonBody, _ := json.Marshal(requestBody)

		mockedDbcon.EXPECT().DBPush(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/postdb", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		r.PostDB(w, req)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		var response service.HTTPModel
		json.Unmarshal(body, &response)
		assert.NotEmpty(t, response.URL)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		requestBody := service.HTTPModel{URL: "invalid-url"}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/postdb", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		r.PostDB(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid URL")
	})

	t.Run("DB Error", func(t *testing.T) {
		requestBody := service.HTTPModel{URL: "https://example.com"}
		jsonBody, _ := json.Marshal(requestBody)

		mockedDbcon.EXPECT().DBPush(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		req := httptest.NewRequest(http.MethodPost, "/postdb", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		r.PostDB(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "db error")
	})
}

func TestGetDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedDbcon := mocks.NewMockDB(ctrl)

	r := &Router{dbCon: mockedDbcon}

	t.Run("Success", func(t *testing.T) {
		requestBody := service.HTTPModel{URL: "fhuOzV3w3b"}
		jsonBody, _ := json.Marshal(requestBody)

		mockedDbcon.EXPECT().DBGet(gomock.Any(), gomock.Any()).Return("https://example.com", nil)
		req := httptest.NewRequest(http.MethodPost, "/postdb", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		r.GetDB(w, req)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		var response service.HTTPModel
		json.Unmarshal(body, &response)
		assert.NotEmpty(t, response.URL)
	})

	t.Run("DB Error", func(t *testing.T) {
		requestBody := service.HTTPModel{URL: "fhuOzV3w3b"}
		jsonBody, _ := json.Marshal(requestBody)

		mockedDbcon.EXPECT().DBGet(gomock.Any(), gomock.Any()).Return("", errors.New("data not found"))

		req := httptest.NewRequest(http.MethodPost, "/postdb", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		r.GetDB(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "data not found")
	})
}

func TestPostCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedCache := mocks.NewMockCacheInterface(ctrl)
	r := &Router{Cache: mockedCache}

	t.Run("Success", func(t *testing.T) {
		requestBody := service.HTTPModel{URL: "https://example.com"}
		jsonBody, _ := json.Marshal(requestBody)

		mockedCache.EXPECT().PushCache(gomock.Any(), gomock.Any()).Return("", nil)

		req := httptest.NewRequest(http.MethodPost, "/postcache", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		r.PostCache(w, req)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		var response service.HTTPModel
		json.Unmarshal(body, &response)
		assert.NotEmpty(t, response.URL)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		requestBody := service.HTTPModel{URL: "invalid-url"}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/postcache", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		r.PostDB(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid URL")
	})
}

func TestGetCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedCache := mocks.NewMockCacheInterface(ctrl)
	r := &Router{Cache: mockedCache}

	t.Run("Success", func(t *testing.T) {
		requestBody := service.HTTPModel{URL: "fhuOzV3w3b"}
		jsonBody, _ := json.Marshal(requestBody)

		mockedCache.EXPECT().GetCache(gomock.Any()).Return("https://example.com", nil)
		req := httptest.NewRequest(http.MethodPost, "/getcache", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		r.GetCache(w, req)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		var response service.HTTPModel
		json.Unmarshal(body, &response)
		assert.NotEmpty(t, response.URL)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/getcache?short_url=fhuOzV3w3b", nil)
		w := httptest.NewRecorder()

		mockedCache.EXPECT().GetCache("fhuOzV3w3b").Return("", errors.New("Short url not found."))

		r.GetCache(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Short url not found.")
	})
}
