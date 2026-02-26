package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var IA_ScrapeBaseURL = "https://archive.org/services/search/v1/scrape?"

const MAX_RESULTS = 5000

type scrapeItems struct {
	Count          int    `json:"count"`
	Cursor         string `json:"cursor"`
	CursorPrevious string `json:"previous"`
	Items          []searchItem
	Total          int `json:"total"`
}

type searchItem struct {
	AddedDate              string      `json:"addeddate"`
	AvgRating_Raw          interface{} `json:"avg_rating"`
	AvgRating              []int
	BTIH                   string      `json:"btih"`
	BackupLocation_Raw     interface{} `json:"backup_location"`
	BackupLocation         []string
	Collection             []string    `json:"collection"`
	CollectionsOrdered     string      `json:"collections_ordered"`
	CurateDate             string      `json:"curatedate"`
	CurateNote_Raw         interface{} `json:"curatenote"`
	CurateNote             []string
	CurateState            string      `json:"curatestate"`
	Curation_Raw           interface{} `json:"curation"`
	Curation               []string
	Curator                string      `json:"curator"`
	Date_Raw               interface{} `json:"date"`
	Date                   []string
	Description            interface{} `json:"description"`
	Downloads              int         `json:"downloads"`
	ExternalMetadataUpdate string      `json:"external_metadata_update"`
	FilesCount             int         `json:"files_count"`
	Format_Raw             interface{} `json:"format"`
	Format                 []string
	//Format              []string    `json:"format"`
	Identifier          string      `json:"identifier"`
	IndexDate           string      `json:"indexdate"`
	ItemSize            int         `json:"item_size"`
	LicenseURL          string      `json:"licenseurl"`
	ListMemberships_Raw interface{} `json:"list_memberships"`
	ListMemberships     []string
	// https://pkg.go.dev/encoding/json#RawMessage
	MatchDateAoustid     string      `json:"match_date_acoustid"`
	MediaType            string      `json:"mediatype"`
	Month                int         `json:"month"`
	NoArchiveTorrent     string      `json:"noarchivetorrent"`
	NumFavorites         int         `json:"num_favorites"`
	OaiUpdateDate_Raw    interface{} `json:"oai_updatedate"`
	OaiUpdateDate        []string
	PrimaryCollection    string      `json:"primary_collection"`
	PublicDate           string      `json:"publicdate"`
	ReportedServer       string      `json:"reported_server"`
	ReviewBody_Raw       interface{} `json:"reviewbody"`
	ReviewBody           []string
	ReviewData           []string    `json:"review_data"`
	Reviewer_Raw         interface{} `json:"reviewer"`
	Reviewer             []string
	ReviewerItemName_Raw interface{} `json:"reviewer_itemname"`
	ReviewerItemname     []string
	Scanner_Raw          interface{} `json:"scanner"`
	Scanner              []string
	Subject_Raw          interface{} `json:"subject"`
	Subject              []string
	SubjectCount         int         `json:"subject_count"`
	Stars_Raw            interface{} `json:"stars"`
	Stars                []int
	Title_Raw            interface{} `json:"title"`
	Title                []string
	Week                 int         `json:"week"`
	Year_Raw             interface{} `json:"year"`
	Year                 []int
}

func ScrapeSearch(query string, maxNumResults int, chunkSize int, c chan []searchItem, client *http.Client, cache *Cache) error {

	if chunkSize < 100 {
		return fmt.Errorf("Requested num results must be > 100")
	}

	if chunkSize > 5000 {
		return fmt.Errorf("ChunkSize number of results requested exceeded")
	}

	go func() {
		cursor := ""

		if chunkSize != 0 {
			query = query + "&count=" + strconv.Itoa(chunkSize)
		}

		count := 0

		for {
			log.Println("-------------------New search---------------------")
			if count >= maxNumResults {
				break
			}

			var tmpItems scrapeItems
			url := IA_ScrapeBaseURL + query

			if cursor != "" {
				url = url + "&cursor=" + cursor
				log.Println("-----------Cursor:", cursor)
			}

			log.Println("search", url)

			err := getUrlJSON(client, url, true, &tmpItems, cursor, cache)
			if err != nil {
				log.Fatal(err)
			}

			err = fixSearchItemFields(tmpItems.Items)
			if err != nil {
				log.Fatal(err)
			}
			c <- tmpItems.Items
			count += tmpItems.Count

			if tmpItems.Cursor == "" {
				break
			}
			cursor = tmpItems.Cursor
		}
		close(c)
	}()
	return nil
}
