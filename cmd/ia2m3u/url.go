package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	//"net"
	"net/http"
	"time"
)

var min time.Duration = 99999999999999
var max time.Duration = 0
var n int64 = 0
var total time.Duration = 0

var backOff time.Duration = 3 * time.Second

// func getUrlJSON(client *http.Client, url string, cachePrefix *string, useCache bool, items *scrapeItems, cursor string, cache *Cache) error {
func getUrlJSON(client *http.Client, url string, useCache bool, alternateKey string, items interface{}, cursor string, cache *Cache) error {

	log.Println("Getting ", url)
	var body []byte

	if useCache {
		var key string
		if alternateKey == "" {
			key = url
		} else {
			key = alternateKey
		}

		if body = cache.GetKey(key); body == nil {
			var res *http.Response
			var err error
			if false {
				res, err = http.Get(url)
				if err != nil {
					fmt.Printf("client: could not create request: %s\n", err)
					log.Fatal()
				}
			} else {

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					fmt.Printf("client: could not create request: %s\n", err)
					log.Fatal()
				}

				startTime := time.Now()
				log.Println("URL start")
				res, err = client.Do(req)
				log.Println("URL end")
				since := time.Since(startTime)

				n++
				total += since
				if since < min {
					min = since
				}

				if since > max {
					max = since
				}

				log.Println(since, min, max, time.Duration(int64(total)/n))

				if since > time.Duration(int64(float64(int64(total))/float64(n)*2.0)) || since > 2*time.Second {
					// Backoff
					backOff = backOff + time.Second + time.Second + time.Second + time.Second
					log.Println(backOff, "getUrlJSON - BACKOFF $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
					log.Println(backOff + time.Second + time.Second + time.Second + time.Second + time.Second)
					time.Sleep(backOff + time.Second + time.Second + time.Second + time.Second + time.Second)
					//log.Println(backOff * (3 * time.Second))
					//time.Sleep(backOff * (3 * time.Second))

				} else {
					if backOff > 0 {
						backOff = backOff - time.Second
					}
				}

				if err != nil {
					fmt.Printf("client: error making http request: %s\n", err)
					log.Fatal()
				}

				if res.StatusCode != 200 {
					body, err = io.ReadAll(res.Body)
					if err == nil {
						log.Println("Error. Response body:")
						log.Println("--------------------------------------------------------------")
						log.Println(string(body))
						log.Println("--------------------------------------------------------------")
					}
					return fmt.Errorf("Failing http code %d (!200)", res.StatusCode)
				}
				if err != nil {
					log.Println("Status code", res.StatusCode)
					log.Println(url)
					log.Println(err)
					log.Fatal(err)
				}
				time.Sleep(100 * time.Millisecond)
			}
			body, err = io.ReadAll(res.Body)
			if err != nil {
				log.Println(err)
				return err
			}
			cache.AddToCache(key, body)

		}

	}

	dec := json.NewDecoder(bytes.NewBuffer(body))

	return dec.Decode(items)
}

func newClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        10,               // Maximum idle connections
		MaxIdleConnsPerHost: 10,               // Maximum idle connections per host
		IdleConnTimeout:     90 * time.Second, // Idle connection timeout
		DisableCompression:  false,            // Enable compression
		DisableKeepAlives:   false,            // Enable keep-alives
		// Dial: (&net.Dialer{
		// 	Timeout:   30 * time.Second,
		// 	KeepAlive: 30 * time.Second,
		// }).Dial,
		// TLSHandshakeTimeout:   15 * time.Second,
		// ResponseHeaderTimeout: 15 * time.Second,
		// ExpectContinueTimeout: 5 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   120 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			log.Println("VIA", len(via), via)
			// Custom redirect handling
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}
}
