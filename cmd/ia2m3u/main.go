package main

import (
	"bufio"
	arg "github.com/alexflint/go-arg"
	ia "github.com/gnewton/iascrape"
	m3u "github.com/k3a/go-m3u"
	"log"
	"math"
	"os"
	"time"
)

// texts, audio, movies, web, image, account, data, collection, software, etree, other
// Default is to create an m3u file in ., referencing URLs for audio
//  -c =  cache load; populates local ID cache; does not generate any m3u file
// -d = directory for m3u file; if missing, dir is "."
// -L = Sound files are local; Downloaded to -d directory
// -q = query
// -i = ID reject list; ascii list, one id per line
// -r = Field reject list; json; form:
//  [
//   "fieldName1": [
//                  "value1"
//                  "value2"
//                 ],
// ]

var SPACE_AND = "%20AND%20"
var AUDIOQUERY = "mediatype%3A(audio)"

type args struct {
	M3UFile       string `arg:"-m,--m3u_file" help:"m3u file" default:"./playlist_ia.m3u"`
	CacheLoad     bool   `arg:"-C,--cache" help:"Run query to load cache; Does not produce any m3u output"`
	Debug         bool   `arg:"-D" help:"Debug mode"`
	Dir           string `arg:"-d,--dir" help:"Directory to write m3u files (and audio if -L)" default:"."`
	IncludeIDList string `arg:"-I,--include" help:"Filename containing one ID per line that is added to the results"`
	LocalAudio    bool   `arg:"-L,--local" help:"m3u references sound files which are downloaded and stored in -d directory"`
	TxtResults    bool   `arg:"-O,--Outputresults" help:"Run query and write results (title, artist, ID) to stdout. Does not produce any m3u output"`
	// Change to queries: Queries  []string `arg:"-q,separate"` see https://github.com/alexflint/go-arg
	Query           string `arg:"-q,--query" help:"The query to run. See https://archive.org/advancedsearch.php for query syntax. Must be URL encoded (i.e. spaces must be %20, equals (\"=\") should be %30, etc. Note %20AND%20mediatype%3A(audio) is appended to query to limit to audio formats"`
	Random          bool   `arg:"-r" help:"Order of audio items in playlist is random"`
	RejectFieldList string `arg:"-F,--rejectfields" help:"Filename containing json map of fieldname1:[value1, value2], fieldname2:[value2, value3]; Fields matching these values are rejected"`
	RejectIDList    string `arg:"-R,--rejectids" help:"Filename containing one ID per line that is rejected"`
	VerifyAudioURL  bool   `arg:"-U" help:"Verifies the URL of the audio file by doing an http HEAD request on the URL"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var args args

	arg.MustParse(&args)

	err, m3uOut := checkArgs(&args)
	if err != nil {
		log.Println(err)
	}

	var file *os.File
	if m3uOut {
		file, err = os.Create(args.M3UFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	itemCache, err := ia.NewCache("cache_item.db")
	if err != nil {
		log.Fatal(err)
	}

	client := ia.NewClient()

	zz := "https://gmail.com"
	err = ia.HeadUrl(client, zz, 5, 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	zz = "https://archive.org/download/78_skye-boat-song_pipes-and-drums-of-h-m-2nd-batt-scots-guards-m-lawson-pipe-major-j-b_gbia3042651b/SKYE BOAT SONG - PIPES AND DRUMS OF H. M. 2nd BATT. SCOTS GUARDS.mp3"
	err = ia.HeadUrl(client, zz, 5, 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	//query := "fields=year,title,collection&q=collection=78%20AND%20mediatype%3Aaudio"
	query := "q=collection%3A78rpm%20AND%20subject%3ABagpipe%20AND%20mediatype%3Aaudio&sorts=btih"
	//query := "fields=title,btih&q=mediatype%3Aaudio&sorts=btih"
	//query := "fields=title,btih&q=title%3Aa*&sorts=btih"

	//query := "fields=title,btih&q=mediatype%3Asoftware&sorts=btih"
	//query := "fields=title,btih&q=mediatype%3Aaudio&sorts=addeddate%20desc"
	//query := "fields=title,btih&q=mediatype%3Atexts&sorts=addeddate&sorts=btih%20desc"
	//query := "fields=title,btih&q=title%3Ab%20AND%20mediatype%3Atexts&sorts=btih&sorts=btih%20desc"
	//query := "fields=title&q=mediatype%3Aaudio"
	//query := "q=mediatype%3A(audio)"
	//query := "q=subject%3A\"Pipe+%26+Drum\""
	//query := "q=title%3A(bagpipe)%20AND%20mediatype%3A(audio)&sorts=title%20desc"

	//query := "fields=*&q=mediatype%3Aaudio&sorts=btih"

	log.Println("ScrapeSearch")

	search := ia.Search{
		Query:      query,
		Client:     client,
		ChunkSize:  5000,
		MaxResults: math.MaxInt64,
		Retries:    5,
	}

	log.Println("Query=", query)

	total, err := search.Total()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(err)
	log.Println("total", total)

	var count int64 = 0
	var m3 *m3u.M3U

	if m3uOut {
		m3 = new(m3u.M3U)
	}

	for {
		results, err := search.Execute()
		if err != nil {
			log.Fatal(err)
		}
		if results == nil {
			break
		}
		log.Println(len(results))

		var item *ia.ItemTopLevelMetadata

		for i := 0; i < len(results); i++ {
			item, err = ia.GetItem(results[i].Identifier, client, itemCache)
			if err != nil {
				log.Fatal(err)
			}
			count = count + 1

			if args.TxtResults {
				outputResults(count, &item.Metadata)
			}

			var records []*m3u.Record
			if m3uOut || args.VerifyAudioURL {
				records = makeM3UEntries(item)
			}
			if m3uOut {
				addAll(m3, records)
			}

			if args.LocalAudio {
				log.Println("LocalAudio - unimplemented")
			}

			if args.VerifyAudioURL {
				for _, rec := range records {
					err := verifyAudio(client, rec.URL)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			if args.CacheLoad {
				// Do nothing
			}

			if args.Debug {
				debug(item)
			}

		}
	}

	if m3uOut {
		w := bufio.NewWriter(file)

		if err := m3.Write(w); err != nil {
			log.Fatal(err)
		}
		w.Flush()
	}
}

var rejectFieldString_ = map[string][]string{
	"creator": []string{
		"BAND OF H.M. SCOTS GUARDS",
		"BAND OF THE SCOTS GUARDS",
		"Band Of H. M. Scots Guards",
		"Band of H.M. Scots Guards",
		"COLDSTREAM",
		"Carole Becker-Douglas",
		"Coldstream",
		"H. M. SCOTS GUARDS BAND",
		"H. Majesty's Scots Guards",
		"His Majesty's Scots Guards Band",
		"Leitung",
		"Mr. R. Everson of the Scots Guards",
		"Regimental",
		"RADERMAN",
		"Gutsul",
		"Gajdos",
		"Full Choir",
		"1st Battalion, The Black Watch (Royal Highland Regiment)",
	},
}

var idList = []string{
	"pipes-of-scotland-glasgow-police-pipe-band-bbc-d.-d.-teoli-jr.-a.-c..",
	"raretunes_364_beating-retreat-edinburgh-castle",
	"bowhill1",
	"1st Battalion, The Black Watch (Royal Highland Regiment)",
	"pipes-of-scotland-glasgow-police-pipe-band-bbc-d.-d.-teoli-jr.-a.-c..",
	//"YPB2010-03-02",
	"Rlpb2012CompetitionsSet",
	"lp_scotland-for-ever_the-royal-scots-greys",
	"lp_champions-of-the-world_the-edinburgh-police-pipe-band",
	"lp_scottish-pipes-and-drums_pipe-major-reids-pipe-band",
	"lp_the-pipes-drums-of-the-1st-battalion-s_1st-battalion-scots-guards",
	"lp_in-concert-en-route_1st-battalion-the-black-watch-royal-highla",
	"lp_scottish-heritage_the-48th-highlanders-of-canada",
	"lp_the-pipes-drums-of-the-1st-battalion-scot_the-pipes-drums-of-the-1st-battalion-scot",
	"lp_the-black-watch_the-band-of-the-black-watch",
	"lp_scottish-soldiers_the-massed-military-bands-of-the-royal",
	"lp_scottish-folk-dances_international-bagpipe-band",
	"lp_here-comes-the-famous-48th_the-48th-highlanders-of-canada",
	"lp_scotlands-pride_the-royal-scots-greys",
	"lp_highland-pageantry_the-regimental-band-and-pipes-and-drums-of",
	"lp_highland-pageantry_the-regimental-band-and-pipes-and-drums-of_0",
	"lp_scots-guards-pipes-and-drums-marches_pipes-and-drums-of-the-scots-guards-joh",
	"lp_marches_pipes-and-drums-of-the-scots-guards-john-s",
	"lp_r-na-bpobair-the-king-of-the-pipers_leo-rowsome",
	"lp_pipes-and-drums-of-the-48th-highlanders_the-48th-highlanders-of-canada",
	"lp_the-scots-guards-on-parade_the-regimental-band-of-the-scots-guards_0",
	"lp_highland-pipes_pipes-and-drums-of-2nd-battalion-scots",
	"lp_kilts-on-parade_st-columcilles-united-gaelic-pipe-band",
}
