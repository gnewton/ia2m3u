package main

import (
	"bufio"
	"cmp"
	"fmt"
	arg "github.com/alexflint/go-arg"
	ia "github.com/gnewton/iascrape"
	m3u "github.com/k3a/go-m3u"
	"log"
	"math"
	"net/url"
	"slices"
	"strconv"
	"strings"
	//"net/url"
	"os"
	//"time"
)

type args struct {
	AccessLevel           AccessLevels `arg:"-A" help:"0 - Unrestricted only;  1 - All;   2 - Restricted only.  Internet archive metadata: access-restricted-item. See https://help.archive.org/help/downloading-a-basic-guide/" default:"0"`
	CacheFile             string       `arg:"-c,--cache-file" help:"Location of item JSON cache file" default:"cache_item.db"`
	CacheLoad             bool         `arg:"-C,--cache" help:"Run query to load cache; Does not produce any m3u output"`
	Debug                 bool         `arg:"-g" help:"Debug mode"`
	Dedup                 bool         `arg:"-D" help:"Deduplicate: Use Year, Title and Creator to decide if is a duplicate."`
	Dir                   string       `arg:"-d,--dir" help:"Directory to write files (and audio if -L)" default:"./"`
	Formats               string       `arg:"-f,--formats" help:"Comma separated list of formats in order of preference. Possible values: 'MP3', 'VBR MP3', '128Kbps MP3', '64Kbps MP3', 'MPEG-4 Audio','Ogg Vorbis', 'WAVE', 'Flac', 'AIFF'. 'VBR MP3' is always appended to supplied list."`
	HTMLResults           bool         `arg:"-H,--simplehtml" help:"Produce simple HTML output to stdout. Does not produce any m3u output"`
	Identifier            string       `arg:"-i,--id" help:"Single archive.org identifier to download"`
	IncludeIDFile         string       `arg:"-I,--include" help:"Filename containing one ID per line that is added to the results"`
	JSONOut               bool         `arg:"-J,--json" help:"Outputs JSON to stdout"`
	Limit                 int64        `arg:"-l,--limit" help:"Limit the results to this number" default:"9223372036854775807"`
	LocalAudio            bool         `arg:"-L,--local" help:"m3u references sound files which are downloaded and stored in -d directory"`
	M3UFile               string       `arg:"-m,--m3u_file" help:"m3u file name. Default is {Metadata.Identifier}.m3u" default:"ia_playlist.m3u"`
	M3UPerHit             bool         `arg:"-U,--unique" help:"Generate an M3U playlist per ID."`
	MinAudioLengthSeconds int64        `arg:"-t,--len" help:"Tracks less than this time in seconds are filtered out" default:"0"`
	MusicFilter           bool         `arg:"-M,--music" help:"Filter out non music. Imperfect."`
	Offset                int64        `arg:"-o,--offset" help:"Offset (skip) this number of results befiesre starting limit count" default:"0"`
	PrintMusicFilterList  bool         `arg:"-z,--music" help:"Filter out non music. Imperfect."`
	Queries               []string     `arg:"-q,--query,separate" help:"The query to run. See https://archive.org/advancedsearch.php for query syntax. Must be URL encoded (i.e. spaces must be %20, equals (\"=\") should be %30, etc. Note %20AND%20mediatype%3A(audio) is appended to query to limit to audio formats"` // Change to queries: Queries  []string `arg:"-q,separate"` see https://github.com/alexflint/go-arg
	Random                bool         `arg:"-r" help:"Order of audio items in playlist is random"`
	RejectFieldsFile      string       `arg:"-F,--rejectfields" help:"Filename containing json map of fieldname1:[value1, value2], fieldname2:[value2, value3]; Fields matching these values are rejected. All strings."`
	RejectIDFile          string       `arg:"-R,--rejectids" help:"Filename containing one ID per line that is rejected"`
	Smallest              bool         `arg:"-s" help:"Select the smallest sized audio file"`
	Strm                  bool         `arg:"-S" help:"Generate one stream per audio file"`
	TitleInLocal          bool         `arg:"-T,--title_in_local" help:"Add the title to the local audio filename. Note can result in very long filenames, some that may be too long for some OSes and/or filestystems."`
	Txtresults            bool         `arg:"-O,--Outputresults" help:"Run query and write results (title, artist, ID) to stdout. Does not produce any m3u output"`
	Verbose               bool         `arg:"-v" help:"Verbose output"`
	VerifyAudioURL        bool         `arg:"-U" help:"Verifies the URL of the audio file by doing an http HEAD request on the URL"`
	YearMapFile           string       `arg:"-Y" help:"File containing 'id year' mappings. ID space YEAR"`
	Years                 []int        `arg:"-y" help:"Limit by year range. Two year values, start end (inclusive). i.e. -y 1980 1990"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	args := new(args)

	arg.MustParse(args)

	if args.PrintMusicFilterList {
		printMusicFilterList()
		return
	}

	//
	m3uOut, err := args.check()
	if err != nil {
		log.Fatal(err)
	}

	m3uOut = true

	if args.Verbose {
		log.Println("Wanted formats: ", args.Formats)
		log.Println("M3UPerHit:", args.M3UPerHit)
	}
	log.Println("M3UPerHit:", args.M3UPerHit)

	var itemCache *ia.Cache
	itemCache = nil

	if args.Identifier == "" {
		itemCache, err = ia.NewCache(args.CacheFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	var rejectFields map[string][]string
	if args.RejectFieldsFile != "" {
		err := loadRejectFieldsFile(args.RejectFieldsFile, &rejectFields, args.Verbose)
		if err != nil {
			log.Fatal(err)
		}

	}

	var idYear map[string]int

	if len(args.YearMapFile) != 0 {
		idYear = loadIDYear(args.YearMapFile)
	}

	client := ia.NewClient()

	var m3 *m3u.M3U
	if m3uOut {
		m3 = new(m3u.M3U)
	}

	offset := args.Offset
	limit := args.Limit

	if args.Verbose {
		log.Println("   Offset", offset)
		log.Println("   Limit", limit)

	}

	var acceptedItems []*ia.ItemTopLevelMetadata

	loadedIDs := make(map[string]struct{})

	if args.IncludeIDFile != "" {
		if args.Verbose {
			log.Println("Loading extras", args.IncludeIDFile)
		}
		err := loadExtraIDs(&acceptedItems, loadedIDs, args, client, itemCache, m3, m3uOut, args.CacheLoad, args.Verbose)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(args.Years) == 2 {
		for i := 0; i < len(args.Queries); i++ {
			args.Queries[i] = args.Queries[i] + " AND date:[" + strconv.Itoa(args.Years[0]) + "-01-01 TO " + strconv.Itoa(args.Years[1]) + "-12-31]"
		}
	}

	if args.Identifier != "" {
		id := args.Identifier
		item, err := ia.GetItem(id, client, itemCache, args.Verbose)
		if err != nil {
			log.Println(err)
		}
		if item == nil {
			log.Fatal(id)
		}
		loadedIDs[id] = struct{}{}

		if len(item.Metadata.Identifier) == 0 {
			log.Println("Missing identifier for results id=", id)
			log.Println(item)
		}

		if err != nil {
			log.Fatal(err)
		}
		if item == nil {
			log.Fatal("Item is nil", id)
		}
		err = processItem(&acceptedItems, item, args, client, itemCache, m3, m3uOut, rejectFields, 0, args.CacheLoad)
		if err != nil {
			log.Fatal(err)
		}
	}

	if args.Verbose {
		log.Println("Query")
	}

	queries := args.Queries

	log.Println(queries)
	log.Println(len(queries))
	log.Println(len(queries) == 0 && args.CacheLoad)
	log.Println(len(queries) > 0 || (len(queries) == 0 && args.CacheLoad))

	if len(queries) > 0 || (len(queries) == 0 && args.CacheLoad) {
		if len(args.Queries) == 0 {
			args.Queries = []string{
				AUDIOQUERY,
			}
		} else {
			for i := 0; i < len(args.Queries); i++ {
				if len(args.Queries[i]) == 0 {
					args.Queries[i] = AUDIOQUERY
				} else {
					args.Queries[i] = args.Queries[i] + SPACE_AND + AUDIOQUERY
				}
			}
		}
		log.Println("MMMMMMMMMMMMMMMMMMMMMMMMMMMMM")

		if len(queries) == 0 && args.CacheLoad {
			queries = []string{
				AUDIOQUERY,
			}
		}

		for _, query := range queries {
			if args.Verbose || args.CacheLoad {
				log.Println("----QQQQQQQQQQQQQQQQQQQQQQQQQQQQQ------------------------------")
				log.Println(query)
				log.Println("---------------------------------------------------------------")
			}

			if args.MusicFilter {
				query = applyMusicFilter(query)
			}

			switch args.AccessLevel {
			case Access_NonRestricted:
				query = query + OnlyNonRestrictedItems_Clause
			case Access_Restricted:
				query = query + OnlyRestrictedItems_Clause
			case Access_All:
			}

			query = "q=" + escapeQuery(query)

			search := ia.Search{
				Query:      query,
				Client:     client,
				ChunkSize:  3001,
				MaxResults: math.MaxInt64,
				Retries:    5,
				Verbose:    args.Verbose,
			}

			if args.Verbose || args.CacheLoad {
				log.Println("Query=", query)
			}

			total, err := search.Total()
			if err != nil {
				log.Fatal(err)
			}
			if args.Verbose {
				log.Println("")
				log.Printf("---- Search total: %d        query: %s\n", total, query)
			}

			var count int = 0
			stop := false
			for {
				if stop {
					break
				}
				results, err := search.Execute()
				if err != nil {
					log.Fatal(err)
				}
				if args.CacheLoad {
					log.Println("Next cursor")
				}
				if results == nil {
					if args.Verbose {
						log.Println("End results")
					}
					break
				}

				for i := 0; i < len(results); i++ {
					if int64(count) > offset {
						id := results[i].Identifier
						if args.Verbose && count%1000 == 0 {
							log.Println(count, "Getting ", results[i].Identifier)
						}
						// Alreaded loaded in this session
						if _, ok := loadedIDs[id]; ok {
							continue
						}

						item, err := ia.GetItem(id, client, itemCache, args.Verbose)
						if err != nil {
							log.Println(err)
						}
						if item == nil {
							continue
						}
						if args.CacheLoad {
							if count%1000 == 0 {
								log.Println(count, id)
							}
						} else {
							loadedIDs[id] = struct{}{}
							if len(item.Metadata.Identifier) == 0 {
								log.Println("Missing identifier for results id=", id)
								log.Println(item)
								continue
							}

							if err != nil {
								log.Fatal(err)
							}
							if item == nil {
								continue
							}
							err = processItem(&acceptedItems, item, args, client, itemCache, m3, m3uOut, rejectFields, count, args.CacheLoad)
							if err != nil {
								log.Fatal(err)
							}
						}
					}

					count++
					if int64(count) > offset+limit {
						stop = true
						break
					}
				}
			}
		}
	}
	var totalTunes = 0
	wantedFormats := makePreferredFormats(args.Formats)
	var dlAudio []DownloadAudio
	adjustYears(idYear, &acceptedItems)

	if args.Dedup {
		acceptedItems = dedup(acceptedItems)
	}

	if m3uOut || args.HTMLResults {
		slices.SortFunc(acceptedItems, itemsByYear) // Order by year
	}

	if m3uOut {
		logV(args.Verbose, "M3U file: "+args.M3UFile)

		m3uFile, err := os.Create(args.M3UFile)
		//file, err := os.Create("playlist_ia.m3u")
		if err != nil {
			panic(err)
		}
		defer m3uFile.Close()

		var thisM3 *m3u.M3U
		for i := 0; i < len(acceptedItems); i++ {
			item := acceptedItems[i]
			if args.M3UPerHit {
				thisM3 = new(m3u.M3U)
				m3uFile, err = os.Create(item.Metadata.Identifier + ".m3u")
				if err != nil {
					panic(err)
				}
			} else {
				thisM3 = m3
			}

			copies := collectCopies(item.Files, args.MinAudioLengthSeconds)
			wantedCopies := findWantedCopies(item.Metadata.Identifier, copies, wantedFormats)

			if len(wantedCopies) > 0 {
				slices.SortFunc(wantedCopies, tunesByTrackOrder) // []*ia.File
				dlAudio = append(dlAudio, makeM3UEntries(item, wantedCopies, thisM3, args.Random, args.LocalAudio)...)
				if args.M3UPerHit {
					w := bufio.NewWriter(m3uFile)
					if err := thisM3.Write(w); err != nil {
						log.Fatal(err)
					}
					w.Flush()
					m3uFile.Close()
				}
			}
		}

		if !args.M3UPerHit {
			w := bufio.NewWriter(m3uFile)
			if err := m3.Write(w); err != nil {
				log.Fatal(err)
			}
			w.Flush()
		}

		if args.LocalAudio {
			err = downloadAudio(dlAudio, args.Verbose)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if args.Verbose {
		log.Println("Output: HTML", args.HTMLResults)
	}

	if args.HTMLResults {
		if args.Verbose {
			log.Println("Output: HTML")
		}
		fmt.Println("<html>")

		for i := 0; i < len(args.Queries); i++ {
			query := args.Queries[i]
			if args.MusicFilter {
				query = applyMusicFilter(query)
			}
			fmt.Printf("<!-- CMD %s     -->\n", os.Args)
			query = "q=" + escapeQuery(query)

			fmt.Printf("<!-- Query %d = [%s]  -->\n", i+1, query)
			equery, err := url.QueryUnescape(query)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("<!-- ### [%s]  -->\n", equery)
		}
		fmt.Printf("\n\n<!-- Offset=%d Limit=%d -->", args.Offset, args.Limit)

		fmt.Println()
		fmt.Println("<body>")
		fmt.Println("<table  style='border-collapse: collapse;' cellpadding='5'>")

		for i := 0; i < len(acceptedItems); i++ {
			//copies := collectCopies(acceptedItems[i].Files)
			copies := collectCopies(acceptedItems[i].Files, args.MinAudioLengthSeconds)
			wantedCopies := findWantedCopies(acceptedItems[i].Metadata.Identifier, copies, wantedFormats)
			if len(wantedCopies) > 0 {
				slices.SortFunc(wantedCopies, tunesByTrackOrder)
				totalTunes += simpleHTML(acceptedItems[i], wantedCopies, args.Verbose)
			}
		}
		fmt.Println("</table>")
		fmt.Println("</body>")
		fmt.Println("</html>")
	}

	if args.Verbose {
		log.Println("Total items:", len(acceptedItems))
		log.Println("Total tunes:", totalTunes)
		hits, misses := ia.CacheStats()
		log.Printf("Cache hit rate: %3.0f%%   %d %d", 100*float64(hits)/float64(hits+misses), hits, misses)
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

func itemsByYear(a, b *ia.ItemTopLevelMetadata) int {
	//if b.Metadata.CanonicalYear == 0 && a.Metadata.CanonicalYear == 0 {
	if b.Metadata.CanonicalYear == a.Metadata.CanonicalYear {
		//if a.Metadata.Source == "Vinyl LP" && a.Metadata.Source != "Vinyl LP" {
		//			return -2
		//}
		//
		// For 78s, will cluster 2 sides
		if strings.Contains(a.Metadata.Identifier, "_gbia") && strings.Contains(b.Metadata.Identifier, "_gbia") {
			aParts := strings.Split(a.Metadata.Identifier, "_gbia")
			bParts := strings.Split(b.Metadata.Identifier, "_gbia")
			return cmp.Compare(aParts[1], bParts[1])
		} else {
			return cmp.Compare(a.Metadata.Titles[0], b.Metadata.Titles[0])
		}
	}

	if b.Metadata.CanonicalYear == a.Metadata.CanonicalYear {
		return cmp.Compare(a.Metadata.Titles[0], b.Metadata.Titles[0])
	}
	return cmp.Compare(b.Metadata.CanonicalYear, a.Metadata.CanonicalYear)
}
