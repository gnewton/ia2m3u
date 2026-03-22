package main

import (
	ia "github.com/gnewton/iascrape"
	m3u "github.com/k3a/go-m3u"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

func m3uOut(title, url string) *m3u.Record {
	log.Println("m3u - not implemented")

	rec := m3u.NewRecord()
	rec.Title = title
	rec.URL = url

	return rec

}

var FileFormats = map[string]struct{}{
	"AIFF":       struct{}{},
	"Flac":       struct{}{},
	"MP3":        struct{}{},
	"Ogg Vorbis": struct{}{},
	"VBR MP3":    struct{}{},
}

var AudioFileBaseUrl = "https://archive.org/download/" // + /{id}/{filename}.mp3

type DownloadAudio struct {
	localFilename string
	remoteUrl     string
}

type FileFormat struct {
	BaseFileName string
	Formats      map[string]struct{}
	File         *ia.File
}

func makeM3UEntries(item *ia.ItemTopLevelMetadata, m3 *m3u.M3U, recMap map[string]*m3u.Record, random bool, local bool, preferredFormats string) []DownloadAudio {
	var download []DownloadAudio

	year := ""
	if len(item.Metadata.Year) > 0 {
		year = " - " + item.Metadata.Year[0]
	}

	title := ""
	if len(item.Metadata.Title) > 0 {
		title = "(" + item.Metadata.Title[0] + ")"
	}

	creator := ""
	if len(item.Metadata.Creator) > 0 {
		creator = item.Metadata.Creator[0] + "(" + year + ") - "
	}

	var formats []string
	if preferredFormats != "" {
		formats = makePreferredFormats(preferredFormats)
	} else {
		formats = []string{"VBR MP3", "MP3", "64Kbps MP3", "128Kbps MP3"}
	}

	log.Println(formats)

	collectedFiles := make(map[string]*FileFormat)

	count := 0
	if len(item.Files) > 0 {
		//log.Println(title, item.Metadata.Identifier)
		for _, file := range item.Files {
			// Flac, WAVE, Ogg Vorbis,
			if isFileFormat(file.Format) {
				log.Println("HELLOOOOOOOOOOOOO")
				collectFile(collectedFiles, &file)
			}
			log.Println("+++++++++++=   ", file.Format, file.Size, file.Name)

			for _, format := range formats {
				if file.Format == format {
					log.Println("Choosing:", file.Format)
					rec := m3u.NewRecord()
					// Tune title
					if len(file.Title) != 0 {
						rec.Title = file.Title
					} else {
						rec.Title = "[Title unknown]"
					}
					rec.Title = creator + title + " -- " + rec.Title
					if local {
						rec.URL = makeLocalAudioURL(item.Metadata.Identifier, file.Name, file.Format, count) // Local
					} else {
						rec.URL = makeRemoteAudioURL(item.Metadata.Identifier, file.Name) // Local
					}

					download = append(download, DownloadAudio{
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
	}

	for _, v := range collectedFiles {
		selected := false
		for thisFormat, _ := range v.Formats {
			log.Println("BBBBBBBB", thisFormat)
			for k := 0; k < len(formats); k++ {
				if formats[k] == thisFormat {
					log.Println("SELECTED", thisFormat)
					selected = true
					break
				}
			}
			if selected {
				break
			}
		}
	}

	return download
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

func makePreferredFormats(pfs string) []string {
	fs := strings.Split(pfs, ",")

	for i := 0; i < len(fs); i++ {
		fs[i] = strings.TrimSpace(fs[i])
	}

	fs = append(fs, "VBR MP3", "MP3", "128Kbps MP3", "64Kbps MP3")
	return fs
}

func isFileFormat(format string) bool {
	_, ok := FileFormats[format]
	return ok
}

func collectFile(collectedFiles map[string]*FileFormat, file *ia.File) {
	baseName := makeBaseName(file.Name)
	log.Println("QQQ", file.Name, baseName)

	var ff *FileFormat
	var ok bool
	if ff, ok = collectedFiles[baseName]; !ok {
		ff = &FileFormat{
			BaseFileName: baseName,
			Formats:      make(map[string]struct{}),
		}
		collectedFiles[baseName] = ff
	}
	ff.Formats[file.Format] = struct{}{}
}

func makeBaseName(f string) string {
	log.Println("A", f)
	f = strings.TrimSuffix(f, filepath.Ext(f))
	log.Println("B", f)
	f = strings.TrimSuffix(f, "_vbr")
	return f
}
