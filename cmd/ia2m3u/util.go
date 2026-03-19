package main

import (
	"errors"
	"fmt"
	ia "github.com/gnewton/iascrape"
	"log"
	"net/http"
	"strings"
	"time"
)

var BadQueryChars = " =\""

func checkArgs(args *args) (error, bool) {
	m3uOut := true
	//Conflicting args
	if args.TxtResults && args.CacheLoad {
		return errors.New("Only one of -O and -C can be true"), false
	}

	if args.TxtResults && args.LocalAudio {
		return errors.New("Only one of -O and -L can be true"), false
	}

	if args.CacheLoad && args.LocalAudio {
		return errors.New("Only one of -C and -L can be true"), false
	}

	// Query
	if strings.ContainsAny(args.Query, BadQueryChars) {
		return errors.New("Query contains unescaped character(s). Cannot contain these characters " + BadQueryChars), false
	}

	if args.TxtResults || args.CacheLoad {
		m3uOut = false
	}

	if len(args.Query) == 0 {
		args.Query = AUDIOQUERY
	} else {
		args.Query = args.Query + SPACE_AND + AUDIOQUERY
	}

	return nil, m3uOut
}

func makeTitle(titles []string) string {
	if len(titles) == 0 {
		return "[Title unknown]"
	}
	return titles[0]
}

func makeURL(item *ia.ItemTopLevelMetadata) string {
	return "HTTP UNKNOWN"
}

func outputResults(count int64, item *ia.ItemMetadata) {
	year := "????"

	if len(item.Year) != 0 && item.Year[0] != "" {
		year = item.Year[0]
	}

	creator := "?"
	if len(item.Creator) != 0 && item.Creator[0] != "" {
		creator = item.Creator[0]
	}

	title := "?"
	if len(item.Title) != 0 && item.Title[0] != "" {
		title = item.Title[0]
	}

	//fmt.Println(count, "Year=", year, " Title=", title, " Creator=", creator, "  ID=", item.Identifier)
	fmt.Printf("%d \t Year=%s \t Title=\"%s\"     Creator=\"%s\"     ID=%s\n", count, year, title, creator, item.Identifier)
	//fmt.Println(year, title, creator)
}

func debug(item *ia.ItemTopLevelMetadata) {
	log.Println(item.Metadata.Identifier)
	if len(item.Files) > 0 {
		for _, file := range item.Files {
			if file.Format == "VBR MP3" {
				log.Println("-----", file.Name, file.Format, file.Title, file.Size)
			}
		}
	}
}

func verifyAudio(client *http.Client, url string) error {
	log.Println("verifyAudio", url)
	return ia.HeadUrl(client, url, 5, 3*time.Second)

}
