package main

import (
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := initializeCache("cache.db")
	if err != nil {
		log.Fatal(err)
	}

	//results, err := search("fields=year,title,collection&q=collection=etree", nil)
	//results, err := search("fields=year,title,collection&q=collection=etree", 200, nil)
	results, totalResults, err := search("fields=year,title,collection&q=collection=78", 7321, 1000, nil)
	//results, totalResults, err := search("fields=*&q=collection=78", 1321, 1000, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Num results returned", len(results))
	log.Println("Total results", totalResults)

	//log.Println(results[0])

	for i, item := range results {
		if i > 2 {
			break
		}
		log.Printf("%+v\n", item)
		//log.Println(i, item)
	}
}
