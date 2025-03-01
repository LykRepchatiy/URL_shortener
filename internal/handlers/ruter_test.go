package handlers

import (
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

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
