package main

import (
	"errors"
	"fmt"
	ia "github.com/gnewton/iascrape"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Need to accept then escape; include "&"
// Use net/url.QueryEscape(s string)
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
	for _, query := range args.Query {
		if strings.ContainsAny(query, BadQueryChars) {
			return errors.New("Query contains unescaped character(s). Cannot contain these characters " + BadQueryChars), false
		}
	}

	if args.TxtResults || args.CacheLoad {
		m3uOut = false
	}

	for i := 0; i < len(args.Query); i++ {
		if len(args.Query[i]) == 0 {
			args.Query[i] = AUDIOQUERY
		} else {
			args.Query[i] = args.Query[i] + SPACE_AND + AUDIOQUERY
		}
	}
	return nil, m3uOut
}

func makeTitle(titles []string) string {
	if len(titles) == 0 {
		return "[Title unknown]"
	}
	return titles[0]
}

//func makeURL(item *ia.ItemTopLevelMetadata) string {
//return "HTTP UNKNOWN"
//}

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

func escapeQuery(q string) string {
	return strings.ReplaceAll(url.PathEscape(q), "=", "%3A")
}

func downloadAudio(downloadUrls []DownloadAudio) error {
	log.Println("=============================================")
	for i := 0; i < len(downloadUrls); i++ {
		log.Println("     ----- Download", downloadUrls[i].remoteUrl, downloadUrls[i].localFilename)
		// Create the file
		out, err := os.Create(downloadUrls[i].localFilename)
		if err != nil {
			return err
		}
		defer out.Close()

		// Get the data
		resp, err := http.Get(downloadUrls[i].remoteUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Check server response
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %s", resp.Status)
		}

		// Writer the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}

	}

	return nil
}
