# Make M3U playlists from Internet Archive audio collection



# CLI

```
$ ./ia2m3u --help
Usage: ia2m3u [--accesslevel ACCESSLEVEL] [--len LEN] [--cache-file CACHE-FILE] [--cache] [--debug] [--dedup] [--dir DIR] [--formats FORMATS] [--simplehtml] [--id ID] [--include INCLUDE] [--limit LIMIT] [--local] [--m3u_file M3U_FILE] [--unique] [--music] [--offset OFFSET] [--query QUERY] [--random] [--rejectfields REJECTFIELDS] [--rejectids REJECTIDS] [--smallest] [--strm] [--title_in_local] [--Outputresults] [--verbose] [--verifyaudiourl] [--yearmapfile YEARMAPFILE] [--years YEARS]

Options:
  --accesslevel ACCESSLEVEL, -A ACCESSLEVEL
                         0 - Unrestricted only;  1 - All;   2 - Restricted only.  Internet archive metadata: access-restricted-item. See https://help.archive.org/help/downloading-a-basic-guide/ [default: 0]
  --len LEN, -t LEN      Tracks less than this time in seconds are filtered out [default: 0]
  --cache-file CACHE-FILE, -c CACHE-FILE
                         Location of item JSON cache file [default: cache_item.db]
  --cache, -C            Run query to load cache; Does not produce any m3u output
  --debug, -g            Debug mode
  --dedup, -D            Deduplicate: Use Year, Title and Creator to decide if is a duplicate.
  --dir DIR, -d DIR      Directory to write files (and audio if -L) [default: ./]
  --formats FORMATS, -f FORMATS
                         Comma separated list of formats in order of preference. Possible values: 'MP3', 'VBR MP3', '128Kbps MP3', '64Kbps MP3', 'MPEG-4 Audio','Ogg Vorbis', 'WAVE', 'Flac', 'AIFF'. 'VBR MP3' is always appended to supplied list.
  --simplehtml, -H       Produce simple HTML output to stdout. Does not produce any m3u output
  --id ID, -i ID         Single archive.org identifier to download
  --include INCLUDE, -I INCLUDE
                         Filename containing one ID per line that is added to the results
  --limit LIMIT, -l LIMIT
                         Limit the results to this number [default: 9223372036854775807]
  --local, -L            m3u references sound files which are downloaded and stored in -d directory
  --m3u_file M3U_FILE, -m M3U_FILE
                         m3u file name. Default is {Metadata.Identifier}.m3u [default: ia_playlist.m3u]
  --unique, -U           Generate an M3U playlist per ID.
  --music, -M            Filter out non music. Imperfect.
  --offset OFFSET, -o OFFSET
                         Offset (skip) this number of results befiesre starting limit count [default: 0]
  --query QUERY, -q QUERY
                         The query to run. See https://archive.org/advancedsearch.php for query syntax. Must be URL encoded (i.e. spaces must be %20, equals ("=") should be %30, etc. Note %20AND%20mediatype%3A(audio) is appended to query to limit to audio formats
  --random, -r           Order of audio items in playlist is random
  --rejectfields REJECTFIELDS, -F REJECTFIELDS
                         Filename containing json map of fieldname1:[value1, value2], fieldname2:[value2, value3]; Fields matching these values are rejected. All strings.
  --rejectids REJECTIDS, -R REJECTIDS
                         Filename containing one ID per line that is rejected
  --smallest, -s         Select the smallest sized audio file
  --strm, -S             Generate one stream per audio file
  --title_in_local, -T   Add the title to the local audio filename. Note can result in very long filenames, some that may be too long for some OSes and/or filestystems.
  --Outputresults, -O    Run query and write results (title, artist, ID) to stdout. Does not produce any m3u output
  --verbose, -v          Verbose output
  --verifyaudiourl, -U   Verifies the URL of the audio file by doing an http HEAD request on the URL
  --yearmapfile YEARMAPFILE, -Y YEARMAPFILE
                         File containing 'id year' mappings. ID space YEAR
  --years YEARS, -y YEARS
                         Limit by year range. Two year values, start end (inclusive). i.e. -y 1980 1990
  --help, -h             display this help and exit
  $
  
```

END
