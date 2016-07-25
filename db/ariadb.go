package ariaDB

import (
	"crypto/md5"
	"encoding/hex"
)

type db struct {
	urls map[string]string
}

func NewDB() *db {
	return &db{
		urls: make(map[string]string),
	}
}

func (d *db) AddLink(url string) string {
	md5OfURL := md5.Sum([]byte(url))
	key := hex.EncodeToString(md5OfURL[0:4])
	d.urls[key] = url
	return key
}

func (d *db) GetLink(key string) string {
	return d.urls[key]
}
