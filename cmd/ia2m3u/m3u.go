package main

import (
	ia "github.com/gnewton/iascrape"
	m3u "github.com/k3a/go-m3u"
	"log"
)

func m3uOut(title, url string) *m3u.Record {
	log.Println("m3u - not implemented")

	rec := m3u.NewRecord()
	rec.Title = title
	rec.URL = url

	return rec

}

var AudioFileBaseUrl = "https://archive.org/download/" // + /{id}/{filename}.mp3

func makeM3UEntries(item *ia.ItemTopLevelMetadata) []*m3u.Record {
	records := make([]*m3u.Record, 0)

	tmp := m3u.NewRecord()
	tmp.Title = "fooTitle"
	tmp.URL = "fooURL"

	year := ""
	if len(item.Metadata.Year) > 0 {
		year = " - " + item.Metadata.Year[0]
	}

	title := ""
	if len(item.Metadata.Title) > 0 {
		title = "(" + item.Metadata.Title[0] + year + ")"
	}

	creator := ""
	if len(item.Metadata.Creator) > 0 {
		creator = item.Metadata.Creator[0] + " - "
	}

	if len(item.Files) > 0 {
		//log.Println(title, item.Metadata.Identifier)
		for _, file := range item.Files {
			if file.Format == "VBR MP3" {

				rec := m3u.NewRecord()
				if len(file.Title) != 0 {
					rec.Title = file.Title
				} else {
					rec.Title = "[Title unknown]"
				}
				rec.Title = creator + title + " -- " + rec.Title
				//rec.URL = AudioFileBaseUrl + item.Metadata.Identifier + "/" + file.Name
				rec.URL = makeAudioURL(item.Metadata.Identifier, file.Name)

				records = append(records, rec)
			}
		}
	}

	return records
}

func makeAudioURL(id, filename string) string {
	return AudioFileBaseUrl + id + "/" + filename
}

func addAll(m3 *m3u.M3U, records []*m3u.Record) {
	for i := 0; i < len(records); i++ {
		m3.Add(records[i])
	}
}
