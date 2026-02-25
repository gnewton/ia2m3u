package main

var baseUrl = "https://archive.org/metadata/"

type ItemTopLevelMetadata struct {
	Created         int64        `json:"created"`
	D1              string       `json:"d1"`
	Date            string       `json:"date"`
	Dir             string       `json:"dir"`
	Files           []File       `json:"files"`
	ItemLastUpdated int          `json:"item_last_updated"`
	ItemSize        int64        `json:"item_size"`
	Metadata        ItemMetadata `json:"metadata"`
	Roles           Role         `json:"roles"`
	Segments        interface{}  `json:"segments"`
}

type ItemMetadata struct {
	CollectionCatalogNumber string      `json:"collection-catalog-number"`
	Creator                 interface{} `json:"creator"`
	Date                    string      `json:"date"`
	DateAsYearInt           int
	Identifier              string      `json:"identifier"`
	Keywords                []string    `json:"keywords"`
	Language                string      `json:"language"`
	Subject                 interface{} `json:"subject"`
	Title                   interface{} `json:"title"`
	Year                    string      `json:"year"`
	MediaType               string      `json:"media_type"`
	Scanner                 string      `json:"scanner"`
	PublicDate              string      `json:"publicdate"`
	AddedDate               string      `json:"addeddate"`
	Uploader                string      `json:"uploader"`
}

type File struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	Title  string `json:"title"`
	Size   string `json:"size"`
}

type Role struct {
	Performer interface{} `json:"performer"`
}

func getItem(id string) {
	//url := baseUrl + id
}
