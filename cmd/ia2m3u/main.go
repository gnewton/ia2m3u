package main

import (
	//"compress/gzip"
	"context"
	ia "github.com/gnewton/iascrape"
	"log"
	"math"
	//"net/http"
	"os"
	"time"
)

// Internet Archive Search api (scrape): https://archive.org/help/aboutsearch.htm

// IA MediaTypes:
// texts, audio, movies, web, image, account, data, collection, software, etree, other
//

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("START")

	itemCache := new(ia.Cache)
	itemCache.KeepForever = true
	err := itemCache.InitializeCache("cache_item.db")
	if err != nil {
		log.Fatal(err)
	}

	client := ia.NewClient()

	//query := "fields=year,title,collection&q=collection=78%20AND%20mediatype%3Aaudio"
	//query := "fields=title,format,btih&q=collection%3A78rpm%20AND%20subject%3ABagpipe%20AND%20mediatype%3Aaudio"
	//query := "fields=title,btih&q=mediatype%3Aaudio&sorts=btih"
	//query := "fields=title,btih&q=title%3Aa*&sorts=btih"

	//query := "fields=title,btih&q=mediatype%3Asoftware&sorts=btih"
	//query := "fields=title,btih&q=mediatype%3Aaudio&sorts=addeddate%20desc"
	//query := "fields=title,btih&q=mediatype%3Atexts&sorts=addeddate&sorts=btih%20desc"
	//query := "fields=title,btih&q=title%3Ab%20AND%20mediatype%3Atexts&sorts=btih&sorts=btih%20desc"
	//query := "fields=title&q=mediatype%3Aaudio"
	query := "q=mediatype%3Aaudio"

	//query := "fields=*&q=mediatype%3Aaudio&sorts=btih"

	log.Println("ScrapeSearch")

	if true {
		scrape := ia.Search{
			Query:      query,
			Client:     client,
			ChunkSize:  5000,
			MaxResults: math.MaxInt64,
		}

		//ctx, cancel := context.WithTimeout(context.Background(), 20000*time.Millisecond)
		//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Millisecond))
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		defer cancel()

		total, err := scrape.Total(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(err)
		log.Println("total", total)

		file, err := os.OpenFile("ids.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		for {
			//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			results, err := scrape.Execute(ctx)
			if err != nil {
				log.Fatal(err)
			}
			if results == nil {
				break
			}
			// log.Println(len(results), counter)
			// counter = counter + len(results)
			for i := 0; i < len(results); i++ {
				// 	log.Println(itemCounter, results[i].Identifier)
				if _, err := file.WriteString(results[i].Identifier + "\n"); err != nil {
					log.Fatal(err)
				}

				ctxItem, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				//_, _, err = ia.GetItem(ctxItem, results[i].Identifier, client, itemCache)
				_, err = ia.GetItem(ctxItem, results[i].Identifier, client, itemCache)
				if err != nil {
					log.Fatal(err)
				} else {
					cancel()
				}
				//itemCounter = itemCounter + 1
			}
		}
	}

	log.Fatal()

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

var rejectFieldString_ = map[string][]string{
	"creator": []string{
		"BAND OF H.M. SCOTS GUARDS",
		"BAND OF THE SCOTS GUARDS",
		"Band Of H. M. Scots Guards",
		"Band of H.M. Scots Guards",
		"COLDSTREAM",
		"Carole Becker-Douglas",
		"Coldstream",
		"H. M. SCOTS GUARDS BAND",
		"H. Majesty's Scots Guards",
		"His Majesty's Scots Guards Band",
		"Leitung",
		"Mr. R. Everson of the Scots Guards",
		"Regimental",
		"RADERMAN",
		"Gutsul",
		"Gajdos",
		"Full Choir",
		"1st Battalion, The Black Watch (Royal Highland Regiment)",
	},
}

var idList = []string{
	"pipes-of-scotland-glasgow-police-pipe-band-bbc-d.-d.-teoli-jr.-a.-c..",
	"raretunes_364_beating-retreat-edinburgh-castle",
	"bowhill1",
	"1st Battalion, The Black Watch (Royal Highland Regiment)",
	"pipes-of-scotland-glasgow-police-pipe-band-bbc-d.-d.-teoli-jr.-a.-c..",
	//"YPB2010-03-02",
	"Rlpb2012CompetitionsSet",
	"lp_scotland-for-ever_the-royal-scots-greys",
	"lp_champions-of-the-world_the-edinburgh-police-pipe-band",
	"lp_scottish-pipes-and-drums_pipe-major-reids-pipe-band",
	"lp_the-pipes-drums-of-the-1st-battalion-s_1st-battalion-scots-guards",
	"lp_in-concert-en-route_1st-battalion-the-black-watch-royal-highla",
	"lp_scottish-heritage_the-48th-highlanders-of-canada",
	"lp_the-pipes-drums-of-the-1st-battalion-scot_the-pipes-drums-of-the-1st-battalion-scot",
	"lp_the-black-watch_the-band-of-the-black-watch",
	"lp_scottish-soldiers_the-massed-military-bands-of-the-royal",
	"lp_scottish-folk-dances_international-bagpipe-band",
	"lp_here-comes-the-famous-48th_the-48th-highlanders-of-canada",
	"lp_scotlands-pride_the-royal-scots-greys",
	"lp_highland-pageantry_the-regimental-band-and-pipes-and-drums-of",
	"lp_highland-pageantry_the-regimental-band-and-pipes-and-drums-of_0",
	"lp_scots-guards-pipes-and-drums-marches_pipes-and-drums-of-the-scots-guards-joh",
	"lp_marches_pipes-and-drums-of-the-scots-guards-john-s",
	"lp_r-na-bpobair-the-king-of-the-pipers_leo-rowsome",
	"lp_pipes-and-drums-of-the-48th-highlanders_the-48th-highlanders-of-canada",
	"lp_the-scots-guards-on-parade_the-regimental-band-of-the-scots-guards_0",
	"lp_highland-pipes_pipes-and-drums-of-2nd-battalion-scots",
	"lp_kilts-on-parade_st-columcilles-united-gaelic-pipe-band",
}
