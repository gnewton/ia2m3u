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

type DownloadAudio struct {
	localFilename string
	remoteUrl     string
}

type FileFormat struct {
	BaseFileName string
	Formats      map[string]struct{}
	File         *ia.File
}

func makeM3UEntries(item *ia.ItemTopLevelMetadata, tunes []*ia.File, m3 *m3u.M3U, recMap map[string]*m3u.Record, random bool, local bool) []DownloadAudio {

	var download []DownloadAudio

	year := strconv.Itoa(item.Metadata.CanonicalYear)

	title := ""
	if len(item.Metadata.Titles) > 0 {
		title = "(" + item.Metadata.Titles[0] + ")"
	}

	creator := ""
	if len(item.Metadata.Creators) > 0 {
		creator = item.Metadata.Creators[0] + "(" + year + ") - "
	}

	for i := 0; i < len(tunes); i++ {
		tuneFile := tunes[i]
		rec := m3u.NewRecord()
		// Tune title
		if len(tuneFile.Title) != 0 {
			rec.Title = tuneFile.Title
		} else {
			rec.Title = "[Title unknown]"
		}
		rec.Title = year + " - " + creator + title + " -- " + rec.Title
		if local {
			rec.URL = makeLocalAudioURL(item.Metadata.Identifier, tuneFile.MD5, tuneFile.Name, tuneFile.Format, i) // Local
		} else {
			rec.URL = makeRemoteAudioURL(item.Metadata.Identifier, tuneFile.Name) // Local
		}

		if local {
			download = append(download, DownloadAudio{
				localFilename: rec.URL,
				remoteUrl:     makeRemoteAudioURL(item.Metadata.Identifier, tuneFile.Name),
			})
		}

		if _, ok := recMap[rec.URL]; !ok {
			recMap[rec.URL] = rec
			if !random {
				m3.Add(rec)
			}
		}

	}

	return download
}

func makeRemoteAudioURL(id, filename string) string {
	return AudioFileBaseUrl + id + "/" + url.PathEscape(filename)
}

func makeLocalAudioURL(id, hash, filename string, format string, n int) string {
	Z := cleanString(strings.TrimSuffix(filename, filepath.Ext(filename)))

	log.Println("345     |", cleanString(filename))
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

	number := ""
	if n < 10 {
		number = "0"
	}
	number = number + strconv.Itoa(n)
	//return id + "_" + number + "_" + hash[len(hash)-8:] + "_" + subtype + "." + suffix
	return id + "__" + number + "_" + Z + "_" + hash[len(hash)-4:] + subtype + "." + suffix
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
	if pfs == "" {
		return []string{VBR_MP3, MP3_128Kbps, MP3_64Kbps, OGG, MPEG_4}
	}

	fs := strings.Split(pfs, ",")

	for i := 0; i < len(fs); i++ {
		fs[i] = strings.TrimSpace(fs[i])
	}

	fs = append(fs, VBR_MP3, MP3_128Kbps, MP3_64Kbps, MPEG_4, OGG)
	return fs
}

func isFileFormat(format string) bool {
	_, ok := FileFormats[format]
	return ok
}

func findTuneCopies(tuneCopies map[string]*FileFormat, coll map[string]map[string]*ia.File, file *ia.File) {
	baseName := makeBaseName(file.Name)

	var ff *FileFormat
	var ok bool
	if ff, ok = tuneCopies[baseName]; !ok {
		ff = &FileFormat{
			BaseFileName: baseName,
			Formats:      make(map[string]struct{}),
			File:         file,
		}
		tuneCopies[baseName] = ff
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
