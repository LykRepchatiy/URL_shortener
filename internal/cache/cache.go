package cache

import (
	"errors"
	"url_shortener/internal/service"
)

type Cache struct {
	Short_origin map[string]string
	Origin_short map[string]string
}

func Init() *Cache {
	var cache Cache
	cache.Origin_short = make(map[string]string)
	cache.Short_origin = make(map[string]string)
	return &cache
}

func (c *Cache) appendToCache(key, value string) {
	c.Short_origin[key] = value
	c.Origin_short[value] = key
}

func (c *Cache) PushCache(key, value string) (string, error) {
	var res string
	
	_, ok := c.Short_origin[key]
	if ok {
		if c.Short_origin[key] == value {
			return key, errors.New("Already exists")
		}
	}

	if ok && c.Short_origin[key] != value {
		tmp := value[:len(value)-1]
		newShort := service.ShortURL(tmp)
		c.appendToCache(newShort, value)
		res = newShort
		return res, nil
	} else {
		c.appendToCache(key, value)
	}

	return res, nil
}

func (c *Cache) GetCache(key string) (string, error) {
	if value, ok := c.Short_origin[key]; ok {
		return value, nil
	}
	return "", errors.New("Short url not found.")
}
