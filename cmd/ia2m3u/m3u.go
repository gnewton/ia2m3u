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
	"128Kbps MP3": struct{}{},
	"64Kbps MP3":  struct{}{},
	"AIFF":        struct{}{},
	"Flac":        struct{}{},
	"MP3":         struct{}{},
	"Ogg Vorbis":  struct{}{},
	"VBR MP3":     struct{}{},
}

type DownloadAudio struct {
	localFilename string
	remoteUrl     string
}

type FileFormat struct {
	BaseFileName string
	Formats      map[string]struct{}
	File         *ia.File
}

func makeM3UEntries(item *ia.ItemTopLevelMetadata, m3 *m3u.M3U, recMap map[string]*m3u.Record, random bool, local bool, preferredFormats string, uniqueAudioFiles map[string]struct{}) []DownloadAudio {

	var download []DownloadAudio

	year := ""
	if len(item.Metadata.Years) > 0 {
		year = " - " + item.Metadata.Years[0]
	}

	title := ""
	if len(item.Metadata.Titles) > 0 {
		title = "(" + item.Metadata.Titles[0] + ")"
	}

	creator := ""
	if len(item.Metadata.Creators) > 0 {
		creator = item.Metadata.Creators[0] + "(" + year + ") - "
	}

	var formats []string
	if preferredFormats != "" {
		formats = makePreferredFormats(preferredFormats)
	} else {
		formats = []string{"VBR MP3", "MP3", "64Kbps MP3", "128Kbps MP3"}
	}

	collectedFiles := make(map[string]*FileFormat)

	// basefilename --> format --> File
	nameFormatFile := make(map[string]map[string]*ia.File)

	/////////////////////
	year = item.Metadata.CanonicalYear

	count := 0
	if len(item.Files) > 0 {
		//log.Println(title, item.Metadata.Identifier)
		for _, file := range item.Files {
			// Flac, WAVE, Ogg Vorbis,
			if isFileFormat(file.Format) {
				if _, ok := uniqueAudioFiles[file.MD5]; ok {
					continue
				} else {
					uniqueAudioFiles[file.MD5] = struct{}{}
				}
				collectFile(collectedFiles, nameFormatFile, &file)
			}
		}
	}

	for _, formatFile := range nameFormatFile {
		selected := false
		for format, file := range formatFile {
			for i := 0; i < len(formats); i++ {
				if format == formats[i] {
					rec := m3u.NewRecord()
					// Tune title
					if len(file.Title) != 0 {
						rec.Title = file.Title
					} else {
						rec.Title = "[Title unknown]"
					}
					rec.Title = year + " - " + creator + title + " -- " + rec.Title
					if local {
						rec.URL = makeLocalAudioURL(item.Metadata.Identifier, file.Name, format, count) // Local
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
	id = strings.TrimRight(id, ".")
	suffix := "mp3"
	subtype := ""

	switch format {
	case "FLAC":
		suffix = "flac"
	case "Ogg Vorbis":
		suffix = "ogg"
	case "AIFF":
		suffix = "aiff"
	case "128Kbps MP3":
		subtype = "_128k"
	case "64Kbps MP3":
		subtype = "_64k"
	case "VBR MP3":
		subtype = "_VBR"
	}

	log.Println("~~~~~~~  Localfile:", suffix, "-", format)
	number := ""
	if n < 10 {
		number = "0"
	}
	number = number + strconv.Itoa(n)
	return id + subtype + "__" + number + "." + suffix
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

func collectFile(collectedFiles map[string]*FileFormat, coll map[string]map[string]*ia.File, file *ia.File) {
	baseName := makeBaseName(file.Name)

	var ff *FileFormat
	var ok bool
	if ff, ok = collectedFiles[baseName]; !ok {
		ff = &FileFormat{
			BaseFileName: baseName,
			Formats:      make(map[string]struct{}),
			File:         file,
		}
		collectedFiles[baseName] = ff
	}
	ff.Formats[file.Format] = struct{}{}
	///////////////////////

	var tuneFormat map[string]*ia.File
	if tuneFormat, ok = coll[baseName]; !ok {
		tuneFormat = make(map[string]*ia.File)
		coll[baseName] = tuneFormat
	}
	tuneFormat[file.Format] = file
}

func makeBaseName(f string) string {
	f = strings.TrimSuffix(f, filepath.Ext(f))
	f = strings.TrimSuffix(f, "_vbr")
	f = strings.TrimSuffix(f, "_64kb")
	f = strings.TrimSuffix(f, "_128kb")
	return f
}
