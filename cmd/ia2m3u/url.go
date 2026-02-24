package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func getUrlJSON(url string, cachePrefix *string, useCache bool, items *searchItems, cursor string) error {

	if db == nil {
		log.Fatal("ERROR: db = nil")
	}

	var body []byte

	if useCache {
		if body = getKey(url); body != nil {
			log.Println("KEY HIT")
		} else {

			log.Println(">>>>>>>>> CACHE MISS")
			res, err := http.Get(url)
			if err != nil {
				log.Println(url)
				log.Println(err)
				log.Fatal(err)
			}

			body, err = io.ReadAll(res.Body)
			if err != nil {
				log.Println(err)
				return err
			}
			addToCache(url, body)
			time.Sleep(2 * time.Second)
		}
	}

	log.Println("Hello")
	dec := json.NewDecoder(bytes.NewBuffer(body))
	err := dec.Decode(items)
	if err != nil {
		log.Println(err)
	}

	return err

}
