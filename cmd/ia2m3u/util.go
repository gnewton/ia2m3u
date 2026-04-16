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
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func (a *args) check() (bool, error) {
	m3uOut := true

	if a.Strm {
		log.Fatal("Not implemented")
	}

	//Conflicting args
	if a.TxtResults && a.CacheLoad {
		return false, errors.New("Only one of -O and -C can be true")
	}

	if a.TxtResults && a.LocalAudio {
		return false, errors.New("Only one of -O and -L can be true")
	}

	if a.CacheLoad && a.LocalAudio {
		return false, errors.New("Only one of -C and -L can be true")
	}

	if len(a.Years) != 2 && len(a.Years) != 0 {
		log.Fatal("Years requries 2 int args: start year end year")
	}

	if len(a.Years) == 2 && a.Years[0] >= a.Years[1] {
		log.Fatal("Start year must be less than end year")
	}

	if a.TxtResults || a.CacheLoad {
		m3uOut = false
	}

	for i := 0; i < len(a.Queries); i++ {
		if len(a.Queries[i]) == 0 {
			a.Queries[i] = AUDIOQUERY
		} else {
			a.Queries[i] = a.Queries[i] + SPACE_AND + AUDIOQUERY
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

func makeTitleCreator(titles, creators []string) (string, string) {

	creator := "?"
	if len(creators) != 0 && creators[0] != "" {
		creator = creators[0]
	}

	title := "?"
	if len(titles) != 0 && titles[0] != "" {
		title = titles[0]
	}
	return title, creator
}

func outputResults(count int, item *ia.ItemMetadata) {
	title, creator := makeTitleCreator(item.Titles, item.Creators)

	fmt.Printf(" %d \t %d \t \"%s\"  -- \"%s\"     ID=%s  Subject=%s  Keywords=%s  Genre=%s  Collection=%s\n", count, item.CanonicalYear, title, creator, item.Identifier, item.Subjects, item.Keywords, item.Genres, item.Collections)
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
	return url.QueryEscape(q)
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
		lfilename := downloadUrls[i].localFilename
		// Create the file
		if checkFileExists(lfilename) {
			if verbose {
				log.Println("Exists", lfilename)
			}
			// If exists and is length 0, delete
			fi, err := os.Stat(lfilename)
			if err != nil {
				return err
			}
			// get the size
			if fi.Size() > 0 {
				continue

			} else {
				if verbose {
					log.Println("Removing zero length file", lfilename)
				}
				err := os.Remove(lfilename)
				if err != nil {
					fmt.Println("Error deleting file:", err)
					return err
				}
			}
		}

		out, err := os.Create(lfilename)
		if err != nil {
			return err
		}
		defer out.Close()

		// Get the data
		resp, err := http.Get(downloadUrls[i].remoteUrl)
		if err != nil {
			log.Println("http get error")
			return err
		}
		defer resp.Body.Close()

		// Check server response
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Bad HTTP status code: %s", resp.Status)
		}

		// Writer the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

type Rejects struct {
	RejectFields map[string][]string `json:"rejects"`
}

func loadRejectFieldsFile(rejectFilename string, rejectFields *map[string][]string, verbose bool) error {
	b, err := os.ReadFile(rejectFilename)
	if err != nil {
		log.Fatalf("Failed to read file: %v\n", err)
	}

	err = json.Unmarshal(b, rejectFields)
	if err != nil {
		log.Println("Error loading rejectFields JSON from", rejectFilename)
		log.Println(err)
	}
	if verbose {
		log.Println("Reject fields: ", *rejectFields)
	}

	return err
}

func processItem(acceptedItems *[]*ia.ItemTopLevelMetadata, item *ia.ItemTopLevelMetadata, args *args, client *http.Client, itemCache *ia.Cache, m3 *m3u.M3U, m3uOut bool, rejectFields map[string][]string, count int, cacheLoad bool) error {

	if len(item.Metadata.Identifier) == 0 {
		log.Println("########################################$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$     Missing identifier????")
		log.Println(item)
		return nil
	}
	if args.Verbose {
		//log.Println("HandleItem: Getting metadata record: ", item.Metadata.Identifier)
	}

	if args.Verbose && count > 0 && count%1000 == 0 {
		log.Println("HandleItem:", count, "-", item.Metadata.Identifier)
	}

	if !cacheLoad{
		if rejectByField(&item.Metadata, rejectFields, args.Verbose) {
			if args.Verbose {
				log.Println("Rejected by field")
			}
			return nil
		}

		*acceptedItems = append(*acceptedItems, item)

		if args.HTMLResults {

			//simpleHTML(count, item)
			return nil
		}

		if args.TxtResults {
			outputResults(count, &item.Metadata)
			return nil
		}
	}
	if args.Debug {
		debug(item)
	}

	return nil
}

func rejectByField(item *ia.ItemMetadata, rejectFields map[string][]string, verbose bool) bool {
	if rejectFields == nil { // Don't rejectcompile
		return false
	}
	mm := ia.MakeMetadataItemFieldMap(item)

	for fieldname, field := range mm {

		if rejectValues, ok := rejectFields[fieldname]; ok {
			for i := 0; i < len(rejectValues); i++ {
				if slices.Contains(*field, rejectValues[i]) {
					if verbose {
						log.Println("----------------- REJECTED", *field, " == ", rejectValues[i])
					}
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

	return lines, nil
}

func loadIDYear(filename string) map[string]int {
	idYear := make(map[string]int)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Error: loadIDYear loading file", filename)
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")

	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			break
		}
		parts := strings.Split(lines[i], " ")
		if len(parts) < 2 {
			log.Println(parts)
			log.Fatal("Error loading IDYear file:", filename, " Line:", i, ":", lines[i])
		}
		year, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal("Error loading IDYear file:", filename, " Line:", i, ":", lines[i])
		}
		id := parts[0]
		idYear[id] = year
	}

	return idYear
}

func loadExtraIDs(acceptedItems *[]*ia.ItemTopLevelMetadata, loadedIDs map[string]struct{}, args *args, client *http.Client, itemCache *ia.Cache, m3 *m3u.M3U, m3uOut bool, cacheLoad bool) error {
	ids, err := loadIncludeIDs(args.IncludeIDFile)
	if err != nil {
		return err
	}
	for i := 0; i < len(ids); i++ {
		id := ids[i]
		if len(id) == 0 || id[0] == '#' {
			continue
		}

		if _, ok := loadedIDs[id]; ok {
			continue
		}

		item, err := ia.GetItem(ids[i], client, itemCache, args.Verbose)
		if err != nil {
			return err
		}

		if item == nil {
			continue
		}
		if !cacheLoad{
			loadedIDs[id] = struct{}{}
		}

		err = processItem(acceptedItems, item, args, client, itemCache, m3, m3uOut, nil, 0, cacheLoad)

		if err != nil {
			return err
		}

	}
	return nil
}

// Replace unknown year with info from -Y file
func adjustYears(idYear map[string]int, acceptedItems *[]*ia.ItemTopLevelMetadata) {
	for i := 0; i < len(*acceptedItems); i++ {
		item := (*acceptedItems)[i]
		id := item.Metadata.Identifier
		if year, ok := idYear[id]; ok {
			item.Metadata.CanonicalYear = year
		}
	}
}

func applyMusicFilter(q string) string {
	for field, values := range musicFilter {
		for _, value := range values {
			q = q + " AND -" + field + ":" + "(" + value + ")"
		}
	}
	return q
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func cleanString(s string) string {
	s = nonAlphanumericRegex.ReplaceAllString(s, " ")
	parts := strings.Fields(s)

	s = ""

	for i := 0; i < len(parts); i++ {
		if i > 0 {
			s = s + "_"
		}
		s = s + parts[i]
	}
	return s
}

// Uses Year + Title + Creator to dedup; simplistic; possible errors
func dedup(items []*ia.ItemTopLevelMetadata) []*ia.ItemTopLevelMetadata {
	uniq := make(map[string]*ia.ItemTopLevelMetadata)

	results := make([]*ia.ItemTopLevelMetadata, 0)

	for i := 0; i < len(items); i++ {
		item := items[i]
		metadata := item.Metadata

		if len(metadata.Creators) == 0 || metadata.CanonicalYear == 0 || len(metadata.Titles[0]) == 0 || len(metadata.Creators[0]) == 0 {
			results = append(results, item)
		} else {
			uniq[strconv.Itoa(metadata.CanonicalYear)+metadata.Titles[0]+metadata.Creators[0]] = item
		}

	}

	for _, item := range uniq {
		results = append(results, item)
	}

	return results
}

func makeStrm(item *ia.ItemTopLevelMetadata, tunes []*ia.File) {
	for i := 0; i < len(tunes); i++ {
		tuneFile := tunes[i]
		makeRemoteAudioURL(item.Metadata.Identifier, tuneFile.Name) // Local
	}

}
