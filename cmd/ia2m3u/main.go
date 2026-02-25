package main

import (
	"log"
)

// Internet Archive Search api (scrape): https://archive.org/help/aboutsearch.htm

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cache := new(Cache)

	err := cache.InitializeCache("cache.db")
	if err != nil {
		log.Fatal(err)
	}
	//query := "fields=year,title,collection&q=collection=78%20AND%20mediatype%3Aaudio"
	query := "fields=title,format.year,license,btih&q=collection%3A78rpm%20AND%20subject%3ABagpipe%20AND%20mediatype%3Aaudio"
	//c, err := searchChannel("fields=*&q=mediatype%3Aaudio&sorts=btih", 20000, 5000, cache)
	c := make(chan []searchItem, 2)
	err = scrapeChannel(query, 20000, 5000, c, cache)

	total := 0

	for results := range c {
		log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++", len(results))
		count := 0
		for i, _ := range results {
			item := results[i]
			//if true || count < 20 {
			//log.Println(item.Format)
			log.Println()
			log.Println(total, "TITLE --- ", item.Title)
			log.Println(" BTIH --- ", item.BTIH)
			log.Println("  IDENTIFIER --- ", item.Identifier)
			//}
			count++
			total++
		}
	}

	log.Println("")
	log.Println("")
	log.Println("")

	log.Fatal()
	//results, err := search("fields=year,title,collection&q=collection=etree", nil)
	//results, err := search("fields=year,title,collection&q=collection=etree", 200, nil)
	//results, totalResults, err := search(query, 20000, 5000, nil, cache)
	//results, totalResults, err := search("fields=year,title,collection,identifier&q=mediatype%3Aaudio", 3000, 5000, nil)
	// results, totalResults, err := search("fields=*&q=mediatype%3Aaudio", 12000, 5000, nil, cache)

	//results, totalResults, err := search("fields=*&q=collection=78", 1321, 1000, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Num results returned", len(results))
	// log.Println("Total results", totalResults)

	// //log.Println(results[0])

	// count := 0
	// for i, _ := range results {
	// 	item := results[i]
	// 	if count > 20 {
	// 		break
	// 	}

	// 	//log.Println(item.Format)
	// 	log.Println("TITLE --- ", item.Title)
	// 	log.Println("FORMAT --- ", item.Format)
	// 	if len(item.Format) > 0 {
	// 		log.Println(item.Date, item.Format, item.Title, item.CurateNote, item.Curation)
	// 	}

	// 	// if len(item.Title) != 0 && item.Title[0] != "" {
	// 	// 	log.Println(i, item.Identifier, item.Year, item.Title)

	// 	// }
	// 	//log.Println(i, item)
	// 	//log.Printf("%d   - %#v\n", i, item)
	// 	count++
	// }
}
