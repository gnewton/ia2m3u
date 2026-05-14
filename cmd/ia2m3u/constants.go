package main

var SPACE_AND = " AND "
var AUDIOQUERY = "mediatype:(audio)"

// Archive.org urls, partial urls and specific filenames
var AudioFileBaseUrl = "https://archive.org/download/" // + /{id}/{filename}.mp3
var ItemBaseUrl = "https://archive.org/details/"
var baseUrl = "https://archive.org/metadata/"

var LPBackcoverImage_Format = "Single Page Processed JP2 ZIP"
var LPBackcoverImage_Suffix = "_jp2.zip"
var LPImagesPHP = "/view_archive.php?archive="
var Thumb = "__ia_thumb.jpg"

// var uilleanSource = "https://commons.wikimedia.org/wiki/File:UilleannPipes.jpg"
const AIFF = "AIFF"
const FLAC = "Flac"

// const MP3 = "MP3"
const MPEG_4 = "MPEG-4 Audio"
const MP3_256Kbps = "256Kbps MP3" // https://archive.org/metadata/AlessandraIntroduzioneailaser
const MP3_128Kbps = "128Kbps MP3"
const MP3_64Kbps = "64Kbps MP3"
const OGG = "Ogg Vorbis"
const VBR_MP3 = "VBR MP3"

const AIFF_SUFFIX = ".aiff"
const FLAC_SUFFIX = ".flac"
const MP3_256Kbps_SUFFIX = "_256kb.mp3"
const MP3_128Kbps_SUFFIX = "_128kb.mp3"
const MP3_64Kbps_SUFFIX = "_64kb.mp3"
const OGG_SUFFIX = ".ogg"
const VBR_MP3_SUFFIX = "_vbr.mp3"
const MPEG_4_SUFFIX = ".m4a"

const NotRestrictedItemsClause = " AND -access-restricted-item:(true) "

var FileFormats = map[string]string{
	MP3_256Kbps: MP3_256Kbps_SUFFIX,
	MP3_128Kbps: MP3_128Kbps_SUFFIX,
	MP3_64Kbps:  MP3_64Kbps_SUFFIX,
	AIFF:        AIFF_SUFFIX,
	FLAC:        FLAC_SUFFIX,
	OGG:         OGG_SUFFIX,
	VBR_MP3:     VBR_MP3_SUFFIX,
	MPEG_4:      MPEG_4_SUFFIX,
}

// AND -subject:(Non-Music) AND -subject:(\"Spoken Word\") AND -subject:(\"Monolog\")  AND -subject:(\"Novelty\") AND -subject:(\"Comedy\") AND -collection:(\"audio_religion\") AND -subject:(\"sample\") AND -subject:(\"interviews\")
var musicFilter map[string][]string = map[string][]string{
	"subject": []string{
		"Comedy",
		"Interview",
		"Monolog",
		"Non-Music",
		"Novelty",
		"Podcast",
		"Poem",
		"Spoken Word",
		"Radio Program",
		"Verse",
		"ballad",
		"interviews",
		"radio interview",
		"sample",
	},
	"collection": []string{
		"audio_religion",
		"samples_only",
		"community",
	},
}
