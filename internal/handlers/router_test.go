package handlers

import (
	"reflect"
	"testing"
	"url_shortener/internal/cache"
	"url_shortener/internal/database"
)

func TestNewRouter(t *testing.T) {
	db := database.DataBase{}
	c := cache.Cache{}
	router := NewRouter(db, &c)

	if router == nil {
		t.Errorf("NewRouter вернул nil, ожидался ненулевой указатель")
	}

	if reflect.TypeOf(router).Kind() != reflect.Ptr {
		t.Errorf("NewRouter вернул не указатель, ожидался указатель на Router")
	}
	if reflect.TypeOf(router).Elem() != reflect.TypeOf(Router{}) {
		t.Errorf("NewRouter вернул неверный тип, ожидался *Router")
	}
}
