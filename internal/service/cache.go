package service

import "errors"

// TO DO second map for cache
// connect with router

func appendToCache(key, value string) {
	cache_map[key] = value
}

func Insert(key, value string) (string, error) {
	var res string

	_, ok := cache_map[key]
	if ok {
		if cache_map[key] == value {
			return key, errors.New("Already exists")
		}
	}

	if ok && cache_map[key] != value {
		tmp := value[:len(value)-1]
		newShort := ShortURL(tmp)
		appendToCache(newShort, value)
		res = newShort
		return res, nil
	} else {
		res = ShortURL(value)
		appendToCache(res, value)
	}

	return res, nil
}

func Select(key string) (string, error) {
	if value, ok := cache_map[key]; ok {
		return value, nil
	}
	return "", errors.New("Short url not found.")
}

