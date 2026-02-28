package main

import (
	"log"
	"net/http"
	"sync"
)

// var ItemBaseUrl = "https://archive.org/metadata/"
var ItemBaseUrl = "http://archive.org/metadata/"

type ItemTopLevelMetadata struct {
	Created         int64        `json:"created"`
	D1              string       `json:"d1"`
	Date            string       `json:"date"`
	Dir             string       `json:"dir"`
	Files           []File       `json:"files"`
	ItemLastUpdated int          `json:"item_last_updated"`
	ItemSize        int64        `json:"item_size"`
	Metadata        ItemMetadata `json:"metadata"`
	Roles           Role         `json:"roles"`
	Segments        []string
	Segments_Raw    interface{} `json:"segments"`
}

type ItemMetadata struct {
	AddedDate                   string      `json:"addeddate"`
	CollectionCatalogNumber_Raw interface{} `json:"collection-catalog-number"`
	CollectionCatalogNumber     []string
	Creator                     []string
	Creator_Raw                 interface{} `json:"creator"`
	Date                        []string
	Date_Raw                    interface{} `json:"date"`
	Identifier                  string      `json:"identifier"`
	Keywords                    string      `json:"keywords"`
	Language                    []string
	Language_Raw                interface{} `json:"language"`
	MediaType                   string      `json:"media_type"`
	PublicDate                  string      `json:"publicdate"`
	Scanner                     []string
	Scanner_Raw                 interface{} `json:"scanner"`
	Subject                     []string
	Subject_Raw                 interface{} `json:"subject"`
	Title                       []string
	Title_Raw                   interface{} `json:"title"`
	Uploader                    string      `json:"uploader"`
	Year                        []string
	Year_Raw                    interface{} `json:"year"`
}

type File struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	Title  string `json:"title"`
	Size   string `json:"size"`
}

type Role struct {
	Performer_Raw interface{} `json:"performer"`
	Performer     []string
}

func getItem(id string, client *http.Client, cache *Cache) *ItemTopLevelMetadata {
	url := ItemBaseUrl + id
	var item ItemTopLevelMetadata

	err := getUrlJSON(client, url, true, id, &item, "", cache)
	if err != nil {
		log.Fatal(err)
	}

	fixItemStrings(&item)

	return &item
}

func getItems(searchItems chan []searchItem, client *http.Client, c chan *ItemTopLevelMetadata, cache *Cache, count *int) {

	log.Println("Starting getItems")
	var wg sync.WaitGroup

	idchan := make(chan string, len(c))

	for i := 0; i < cap(c); i++ {
		log.Println("go itemGetter(idchan, c)")
		wg.Add(1)
		go itemGetter(i, &wg, idchan, c, client, cache)
	}

	log.Println("Starting getItems: loop")
	for searchResults := range searchItems {
		for i, _ := range searchResults {
			id := searchResults[i].Identifier
			idchan <- id
			*count++
		}
	}

	log.Println("CLOSING idchan CHANNEL")
	close(idchan)

	wg.Wait()
	close(c)
}

func itemGetter(i int, wg *sync.WaitGroup, ids chan string, items chan *ItemTopLevelMetadata, client *http.Client, cache *Cache) {
	defer wg.Done()
	log.Println("itemGetter START", i)
	for id := range ids {
		//tmp := new(ItemTopLevelMetadata)
		//tmp.Metadata.Identifier = id
		//log.Println(i, id)
		tmd := getItem(id, client, cache)

		items <- tmd
	}
	//log.Println("itemGetter END", i)
}

func fixItemStrings(tm *ItemTopLevelMetadata) error {

	sf := []StringFields{
		{&tm.Segments, tm.Segments_Raw},
		{&tm.Metadata.Subject, tm.Metadata.Subject_Raw},
		{&tm.Metadata.Creator, tm.Metadata.Creator_Raw},
		{&tm.Metadata.Title, tm.Metadata.Title_Raw},
		{&tm.Metadata.Year, tm.Metadata.Year_Raw},
		{&tm.Metadata.Language, tm.Metadata.Language_Raw},
		{&tm.Metadata.Scanner, tm.Metadata.Scanner_Raw},
		{&tm.Metadata.Date, tm.Metadata.Date_Raw},
		{&tm.Metadata.CollectionCatalogNumber, tm.Metadata.CollectionCatalogNumber_Raw},

		{&tm.Roles.Performer, tm.Roles.Performer_Raw},
	}

	//cleanInts(ints, intsRaw)
	err := fixStrings(sf)
	if err != nil {
		//log.Printf("------ %#v\n", item)
		log.Println(err)
		//log.Fatal(item.Identifier)
	}
	return nil
}
