package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// Search api (scrape): https://archive.org/help/aboutsearch.htm

var IA_SearchBaseURL = "https://archive.org/services/search/v1/scrape?"

const MAX_RESULTS = 5000

type searchItems struct {
	Items  []searchItem
	Cursor string `json:"cursor"`
	Count  int    `json:"count"`
	Total  int    `json:"total"`
}

type searchItem struct {
	AddedDate          string   `json:"addeddate"`
	BTIH               string   `json:"btih"`
	CurateDate         string   `json:"curatedate"`
	CurateNote         string   `json:"curatenote"`
	Curation           string   `json:"curation"`
	Curator            string   `json:"curator"`
	Date               string   `json:"date"`
	Description        string   `json:"description"`
	Downloads          int      `json:"downloads"`
	FilesCount         int      `json:"files_count"`
	Format             []string `json:"format"`
	Identifier         string   `json:"identifier"`
	IndexDate          string   `json:"indexdate"`
	ItemSize           int      `json:"item_size"`
	ListMemberships    string   `json:"list_memberships"`
	MediaType          string   `json:"mediatype"`
	Month              int      `json:"month"`
	NumFavorites       string   `json:"num_favorites"`
	OaiUpdateDate      []string `json:"oai_updatedate"`
	PrimaryCollectione string   `json:"primary_collection"`
	PublicDate         string   `json:"publicdate"`
	ReportedServer     string   `json:"reported_server"`
	Scanner            string   `json:"scanner"`
	Subject            string   `json:"subject"`
	SubjectCount       int      `json:"subject_count"`
	Week               int      `json:"week"`
	Year               int      `json:"year"`
	Collection         []string `json:"collection"`
}

func search(query string, wantedNum int, chunkSize int, items []searchItem) ([]searchItem, int, error) {
	if chunkSize < 0 {
		return nil, 0, fmt.Errorf("Num results cannot be < 0")
	}

	if chunkSize > 0 && chunkSize < 100 {
		return nil, 0, fmt.Errorf("Requested num results must be > 100")
	}

	// if chunkSize > 5000 {
	// 	return nil, fmt.Errorf("ChunkSize number of results requested exceeded %d > %d", chunkSize, CHUNKSIZE_RESULTS)
	// }

	if items == nil {
		tmp := new([]searchItem)
		items = *tmp
	}

	cursor := ""

	if chunkSize != 0 {
		query = query + "&count=" + strconv.Itoa(chunkSize)
	}

	count := 0
	var totalResults int

	for {
		if count >= wantedNum {
			break
		}

		var tmpItems searchItems
		url := IA_SearchBaseURL + query

		if cursor != "" {
			url = url + "&cursor=" + cursor
			log.Println("Cursor:", cursor)
		}

		log.Println("search", url)

		// io.Reader
		body, err := getUrlBody(url, nil, false)
		if err != nil {
			return nil, 0, err
		}

		dec := json.NewDecoder(body)

		dec.Decode(&tmpItems)
		count += tmpItems.Count
		totalResults = tmpItems.Total

		log.Println(len(tmpItems.Items))

		items = append(items, tmpItems.Items...)

		if tmpItems.Cursor == "" {
			break
		}
		cursor = tmpItems.Cursor
	}
	return items, totalResults, nil
}
