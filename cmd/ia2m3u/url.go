package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

func getUrlBody(url string, cachePrefix *string, useCache bool) (io.Reader, error) {

	if db == nil {
		log.Fatal("ERROR: db = nil")
	}

	if useCache {
		if value, ok := getKey(url); ok {
			return value, nil
		}
	}

	res, err := http.Get(url)
	if err != nil {
		log.Println(url)
		log.Println(err)
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if useCache {
		addToCache(url, b)
	}

	time.Sleep(2 * time.Second)

	return bytes.NewReader(b), err

}
