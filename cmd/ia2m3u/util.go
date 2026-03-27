package main

import (
	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	ia "github.com/gnewton/iascrape"
	m3u "github.com/k3a/go-m3u"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

var SPACE_AND = "%20AND%20"
var AUDIOQUERY = "mediatype%3A(audio)"

func checkArgs(args *args) (bool, error) {
	m3uOut := true
	//Conflicting args
	if args.TxtResults && args.CacheLoad {
		return false, errors.New("Only one of -O and -C can be true")
	}

	if args.TxtResults && args.LocalAudio {
		return false, errors.New("Only one of -O and -L can be true")
	}

	if args.CacheLoad && args.LocalAudio {
		return false, errors.New("Only one of -C and -L can be true")
	}

	if len(args.Years) != 2 && len(args.Years) != 0 {
		log.Fatal("Years requries 2 int args: start year end year")
	}

	if len(args.Years) == 2 && args.Years[0] >= args.Years[1] {
		log.Fatal("Start year must be less than end year")
	}

	if args.TxtResults || args.CacheLoad {
		m3uOut = false
	}

	for i := 0; i < len(args.Queries); i++ {
		if len(args.Queries[i]) == 0 {
			args.Queries[i] = AUDIOQUERY
		} else {
			args.Queries[i] = args.Queries[i] + SPACE_AND + AUDIOQUERY
		}
	}

	return m3uOut, nil
}

func makeTitle(titles []string) string {
	if len(titles) == 0 {
		return "[Title unknown]"
	}
	return titles[0]
}

func outputResults(count int64, item *ia.ItemMetadata) {
	year := "????"

	if len(item.Years) != 0 && item.Years[0] != "" {
		year = item.Years[0]
	}

	year = item.CanonicalYear

	creator := "?"
	if len(item.Creators) != 0 && item.Creators[0] != "" {
		creator = item.Creators[0]
	}

	title := "?"
	if len(item.Titles) != 0 && item.Titles[0] != "" {
		title = item.Titles[0]
	}

	fmt.Printf(" %d \t %s \t \"%s\"  -- \"%s\"     ID=%s  Subject=%s  Keywords=%s  Genre=%s  Collection=%s\n", count, year, title, creator, item.Identifier, item.Subjects, item.Keywords, item.Genres, item.Collections)
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

func verifyAudio(client *http.Client, url string, verbose bool) error {
	if verbose {
		log.Println("VerifyAudio: Getting HEAD of URL:", url)
	}
	return ia.HeadUrl(client, url, 5, 3*time.Second)

}

func escapeQuery(q string) string {
	return strings.ReplaceAll(url.PathEscape(q), "=", "%3A")
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func downloadAudio(downloadUrls []DownloadAudio, verbose bool) error {

	for i := 0; i < len(downloadUrls); i++ {
		if verbose {
			log.Printf("  ----- Download URL: %s   to local file: %s\n", downloadUrls[i].remoteUrl, downloadUrls[i].localFilename)
		}
		// Create the file
		if checkFileExists(downloadUrls[i].localFilename) {
			if verbose {
				log.Println("Exists")
			}
			continue
		}
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

type Rejects struct {
	RejectFields map[string][]string `json:"rejects"`
}

func loadRejectFieldsFile(rejectFilename string, rejectFields *map[string][]string) error {
	b, err := os.ReadFile(rejectFilename)
	if err != nil {
		log.Fatalf("Failed to read file: %v\n", err)
	}

	err = json.Unmarshal(b, rejectFields)
	log.Println(*rejectFields)

	return err
}

func handleItem(item *ia.ItemTopLevelMetadata, args *args, client *http.Client, itemCache *ia.Cache, recMap map[string]*m3u.Record, m3 *m3u.M3U, m3uOut bool, rejectFields map[string][]string, uniqueAudioFiles map[string]struct{}, count int64) error {
	if args.Verbose {
		log.Println("Getting metadata record: ", item.Metadata.Identifier)
	}

	if rejectByField(&item.Metadata, rejectFields) {
		if args.Verbose {
			log.Println("Rejected by field")
		}
		return nil
	}

	if args.TxtResults {
		outputResults(count, &item.Metadata)
		return nil
	}

	var downloadUrls []DownloadAudio
	if m3uOut || args.VerifyAudioURL {
		downloadUrls = makeM3UEntries(item, m3, recMap, args.Random, args.LocalAudio, args.Formats, uniqueAudioFiles)
	}

	if args.LocalAudio {
		downloadAudio(downloadUrls, args.Verbose)
	}

	if args.VerifyAudioURL {
		log.Println("******************************************", len(downloadUrls))
		for _, url := range downloadUrls {
			err := verifyAudio(client, url.remoteUrl, args.Verbose)
			if err != nil {
				return err
			}
		}
	}
	if args.CacheLoad {
		// Do nothing
	}

	if args.Debug {
		debug(item)
	}

	return nil
}

func rejectByField(item *ia.ItemMetadata, rejectFields map[string][]string) bool {
	if rejectFields == nil { // Don't rejectcompile
		return false
	}
	mm := ia.MakeMetadataItemFieldMap(item)

	for fieldname, field := range mm {
		log.Println(fieldname, field, len(*field))
		if rejectValues, ok := rejectFields[fieldname]; ok {
			for i := 0; i < len(rejectValues); i++ {
				if slices.Contains(*field, rejectValues[i]) {
					log.Println("----------------- REJECTED", *field, " == ", rejectValues[i])
					return true
				}
			}
		}
	}

	return false
}

func loadIncludeIDs(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")

	for i := 0; i < len(lines); i++ {
		log.Println(i, ">>>>>>>>>>>>>>>>  ", lines[i])
	}

	return lines, nil
}

func loadExtraIDs(args *args, client *http.Client, itemCache *ia.Cache, recMap map[string]*m3u.Record, m3 *m3u.M3U, m3uOut bool, uniqueAudioFiles map[string]struct{}) error {
	ids, err := loadIncludeIDs(args.IncludeIDFile)
	if err != nil {
		return err
	}
	for i := 0; i < len(ids); i++ {

		if len(ids[i]) == 0 {
			continue
		}
		item, err := ia.GetItem(ids[i], client, itemCache, args.Verbose)
		if err != nil {
			return err
		}

		err = handleItem(item, args, client, itemCache, recMap, m3, m3uOut, nil, uniqueAudioFiles, 0)

		if err != nil {
			return err
		}

	}
	return nil
}
