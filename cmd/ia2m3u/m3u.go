package main

import (
	ia "github.com/gnewton/iascrape"
	m3u "github.com/k3a/go-m3u"
	"log"
	"net/url"
	"strconv"
)

func m3uOut(title, url string) *m3u.Record {
	log.Println("m3u - not implemented")

	rec := m3u.NewRecord()
	rec.Title = title
	rec.URL = url

	return rec

}

var AudioFileBaseUrl = "https://archive.org/download/" // + /{id}/{filename}.mp3

type DownloadAudio struct {
	localFilename string
	remoteUrl     string
}

func makeM3UEntries(item *ia.ItemTopLevelMetadata, m3 *m3u.M3U, recMap map[string]*m3u.Record, random bool, local bool) []DownloadAudio {
	var dla []DownloadAudio

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

	count := 0
	if len(item.Files) > 0 {
		//log.Println(title, item.Metadata.Identifier)
		for _, file := range item.Files {
			// Flac, WAVE, Ogg Vorbis,
			if file.Format == "Flac" || file.Format == "WAVE" || file.Format == "Ogg Vorbis" || file.Format == "AIFF" || file.Format == "MP3" || file.Format == "VBR MP3" {
				log.Println("+++++++++++=   ", file.Format, file.Size, file.Name)
			}
			if file.Format == "VBR MP3" {
				rec := m3u.NewRecord()
				if len(file.Title) != 0 {
					rec.Title = file.Title
				} else {
					rec.Title = "[Title unknown]"
				}
				rec.Title = creator + title + " -- " + rec.Title
				//rec.URL = AudioFileBaseUrl + item.Metadata.Identifier + "/" + file.Name
				rec.URL = makeLocalAudioURL(item.Metadata.Identifier, file.Name, file.Format, count) // Local

				dla = append(dla, DownloadAudio{
					localFilename: rec.URL,
					remoteUrl:     makeRemoteAudioURL(item.Metadata.Identifier, file.Name),
				})

				if _, ok := recMap[rec.URL]; !ok {
					recMap[rec.URL] = rec
					if !random {
						m3.Add(rec)
					}
				}
				count++
			}
		}
	}

	if count > 0 {
		log.Println("#Items", count)
	}

	return dla
}

func makeRemoteAudioURL(id, filename string) string {
	return AudioFileBaseUrl + id + "/" + url.PathEscape(filename)
}

func makeLocalAudioURL(id, filename string, format string, n int) string {
	var suffix string

	switch format {
	case "MP3":
	case "VBR MP3":
		suffix = "mp3"
	}
	return id + "__" + strconv.Itoa(n) + "." + suffix
}

func addAll(m3 *m3u.M3U, records []*m3u.Record) {
	for i := 0; i < len(records); i++ {
		m3.Add(records[i])
	}
}

func randomizeAudio(m3 *m3u.M3U, recMap map[string]*m3u.Record) {
	for _, value := range recMap {
		m3.Add(value)
	}

}
