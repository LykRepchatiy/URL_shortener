package cache

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	cache := Init()

	if cache == nil {
		t.Errorf("Init() вернул nil, ожидался ненулевой указатель")
	}

	if cache.Origin_short == nil {
		t.Errorf("Origin_short не был инициализирован")
	} else if len(cache.Origin_short) != 0 {
		t.Errorf("Origin_short должен быть пустым, но содержит %d элементов", len(cache.Origin_short))
	}

	if cache.Short_origin == nil {
		t.Errorf("Short_origin не был инициализирован")
	} else if len(cache.Short_origin) != 0 {
		t.Errorf("Short_origin должен быть пустым, но содержит %d элементов", len(cache.Short_origin))
	}

	if reflect.TypeOf(cache.Origin_short).Kind() != reflect.Map {
		t.Errorf("Тип Origin_short не является map")
	}
	if reflect.TypeOf(cache.Short_origin).Kind() != reflect.Map {
		t.Errorf("Тип Short_origin не является map")
	}
}

func TestAppendToCache(t *testing.T) {
    cache := &Cache{
        Origin_short: make(map[string]string),
        Short_origin: make(map[string]string),
    }

    cache.AppendToCache("short1", "origin1")
    if cache.Short_origin["short1"] != "origin1" {
        t.Errorf("Short_origin[\"short1\"] = %s; ожидалось \"origin1\"", cache.Short_origin["short1"])
    }
    if cache.Origin_short["origin1"] != "short1" {
        t.Errorf("Origin_short[\"origin1\"] = %s; ожидалось \"short1\"", cache.Origin_short["origin1"])
    }
}

func TestPushCache(t *testing.T) {
    cache := &Cache{
        Origin_short: make(map[string]string),
        Short_origin: make(map[string]string),
    }

    cache.Short_origin["existingKey"] = "existingValue"
    result, err := cache.PushCache("existingKey", "existingValue")
    if result != "existingKey" || err == nil || err.Error() != "Already exists" {
        t.Errorf("Тест 1: ожидалось ('existingKey', 'Already exists'), получено ('%s', '%v')", result, err)
    }

    cache.Short_origin["existingKey"] = "oldValue"
    result, err = cache.PushCache("existingKey", "newValue")
    expectedShort := "qG60SzSNGn"
    if result != expectedShort || err != nil {
        t.Errorf("Тест 2: ожидалось ('%s', nil), получено ('%s', '%v')", expectedShort, result, err)
    }
    if cache.Short_origin[expectedShort] != "newValue" {
        t.Errorf("Тест 2: новое значение не добавлено в Short_origin")
    }
}

func TestGetCache(t *testing.T) {
    cache := &Cache{
        Origin_short: make(map[string]string),
        Short_origin: make(map[string]string),
    }

    cache.Short_origin["short1"] = "origin1"
    cache.Short_origin["short2"] = "origin2"

    result, err := cache.GetCache("short1")
    if result != "origin1" || err != nil {
        t.Errorf("Тест 1: ожидалось ('origin1', nil), получено ('%s', '%v')", result, err)
    }

    result, err = cache.GetCache("nonexistentKey")
    if result != "" || err == nil || err.Error() != "Short url not found." {
        t.Errorf("Тест 2: ожидалось ('', 'Short url not found.'), получено ('%s', '%v')", result, err)
    }
}