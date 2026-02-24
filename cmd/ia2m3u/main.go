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

	c, err := search3("fields=*&q=mediatype%3Aaudio", 6000, 5000)

	for results := range c {
		log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++", len(results))
		count := 0
		for i, _ := range results {
			item := results[i]
			if count > 20 {
				break
			}

			//log.Println(item.Format)
			log.Println("TITLE --- ", item.Title)
			count++
		}
	}

	log.Fatal()

	//results, err := search("fields=year,title,collection&q=collection=etree", nil)
	//results, err := search("fields=year,title,collection&q=collection=etree", 200, nil)
	//results, totalResults, err := search2("fields=year,title,collection&q=collection=78%20AND%20mediatype%3Aaudio", 101, 1000, nil)
	//results, totalResults, err := search("fields=year,title,collection,identifier&q=mediatype%3Aaudio", 3000, 5000, nil)
	results, totalResults, err := search("fields=*&q=mediatype%3Aaudio", 6000, 5000, nil)

	//results, totalResults, err := search("fields=*&q=collection=78", 1321, 1000, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Num results returned", len(results))
	log.Println("Total results", totalResults)

	//log.Println(results[0])

	count := 0
	for i, _ := range results {
		item := results[i]
		if count > 20 {
			break
		}

		//log.Println(item.Format)
		log.Println("TITLE --- ", item.Title)
		log.Println("FORMAT --- ", item.Format)
		if len(item.Format) > 0 {
			log.Println(item.Date, item.Format, item.Title, item.CurateNote, item.Curation)
		}

		// if len(item.Title) != 0 && item.Title[0] != "" {
		// 	log.Println(i, item.Identifier, item.Year, item.Title)

		// }
		//log.Println(i, item)
		//log.Printf("%d   - %#v\n", i, item)
		count++
	}
}
