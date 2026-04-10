package main

import (
	"cmp"
	"fmt"
	ia "github.com/gnewton/iascrape"
	"net/url"
	//`"slices"
	"strconv"
	"strings"
)

var CREATOR_SEARCH_PREFIX = "https://archive.org/search?query=creator%3A%22"

func simpleHTML(item *ia.ItemTopLevelMetadata, wantedCopies []*ia.File, verbose bool) int {

	meta := item.Metadata

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("<tr bgcolor='eeeeee'>")

	fmt.Printf("<td rowspan='%d' valign='top' align='right' width='5%%'>\n", len(wantedCopies)+1)

	thumb := "https://" + item.D1 + item.Dir + "/" + Thumb

	if has, jp2f := hasJP2ZipFile(item.Files, item.Metadata.Identifier); !has {
		// LP FRONT cover image
		fmt.Println("<a href=\"https://" + item.D1 + item.Dir + "/" + meta.Identifier + "_itemimage.jpg\">")
		fmt.Println("<img width=160 align='left'   style='float: left;'    src='" + thumb + "'>")
		fmt.Println("</a>")
	} else {
		// LP FRONT cover image
		fmt.Println("<a href=\"" + makeJP2ImageUrl(jp2f, item, "0") + "\">")
		fmt.Println("<img width=160 align='left'   style='float: left;'    src='" + thumb + "'>")
		fmt.Println("</a>")

		// LP BACK cover image
		jp2ImageUrl := makeJP2ImageUrl(jp2f, item, "1")
		fmt.Println("<a href=\"" + jp2ImageUrl + "\">")
		fmt.Println("<img width=160 align='left'   style='float: left;'    src='" + jp2ImageUrl + "'>")
		fmt.Println("</a>")
	}
	fmt.Println("</td>")

	// COLUMN
	fmt.Println("")
	fmt.Println("<td valign='top' colspan='3'  bgcolor='eeeeee'>")
	fmt.Println("<h2>")

	var year string
	if meta.CanonicalYear == 0 {
		year = "[Unknown year]"
	} else {
		year = strconv.Itoa(meta.CanonicalYear)
	}

	title, creator := makeTitleCreator(meta.Titles, meta.Creators)

	fmt.Printf("%s <a href=\"https://archive.org/details/%s\">%s</a>\n", year, meta.Identifier, title)

	if creator != "?" {
		fmt.Printf(" - <i><a href=\"%s\">%s</a></i>\n", CREATOR_SEARCH_PREFIX+creator+"%22", creator)
	} else {
		fmt.Printf(" - <i>%s</i>\n", creator)
	}

	if verbose {
		fmt.Printf(" <a href=\"https://archive.org/metadata/%s\">JSON</a>\n", meta.Identifier)
		fmt.Println(meta.Collections)
	}
	fmt.Println("</h2>")
	fmt.Println("</td>")
	fmt.Println("</tr>")

	if len(item.Files) > 0 {
		writeAudioFiles(wantedCopies, meta.Identifier, verbose)
	}

	fmt.Println("<tr   bgcolor='eeeeee'>  <td colspan='4'> <hr> </td> </tr>")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")

	return len(wantedCopies)
}

func tunesByTrackOrder(a, b *ia.File) int {
	return cmp.Compare(a.TrackOrder, b.TrackOrder)
}

func writeAudioFiles(wantedCopies []*ia.File, id string, verbose bool) {

	filenameTitle := make(map[string]string)
	n := 0
	for i := 0; i < len(wantedCopies); i++ {

		f := wantedCopies[i]
		if f == nil {
			continue
		}
		if _, ok := FileFormats[f.Format]; ok { // REDUNDEDNT; REMOVE; FIXXX
			// ROW
			n++
			fmt.Println("")
			fmt.Println("<tr valign='top'>")

			// COLUMN Track#

			if i%2 == 0 {
				fmt.Println("<td  width='2%'  valign='top' align='right'>")
			} else {
				fmt.Println("<td width='2%' bgcolor='eeeeee' valign='top'  align='right' >")
			}

			if len(wantedCopies) > 1 {
				fmt.Printf("%d.  &nbsp; \n", n)
			}
			if len(wantedCopies) == 1 {
				fmt.Printf("%s\n", f.MD5[len(f.MD5)-4:])
			}

			fmt.Println("</td>")

			// COLUMN tune title
			if i%2 == 0 {
				fmt.Println("<td width='35%'  valign='top'>")
			} else {
				fmt.Println("<td width='35%' bgcolor='eeeeee'  valign='top'>")
			}
			//fmt.Printf("<a href=\"%s\">%s</a>", makeRemoteAudioURL(id, f.Name), makeFileTitle(f.Title, f.Name, f.Original, filenameTitle))
			fmt.Printf("%s\n", makeFileTitle(f.Title, f.Name, f.Original, filenameTitle))

			if verbose {
				//fmt.Println(f.TrackOrder, f.Format)
			}
			fmt.Println("</td>")

			// COLUMN - Audio Player
			if i%2 == 0 {
				fmt.Println("<td  valign='top'>")
			} else {
				fmt.Println("<td bgcolor='eeeeee'  valign='top'>")
			}

			fmt.Println("      <audio controls>")
			//fmt.Print("        <source preload='none' src=\"")
			fmt.Print("        <source src=\"")
			fmt.Print(AudioFileBaseUrl + url.PathEscape(id) + "/" + url.PathEscape(f.Name))
			fmt.Print("\"")
			fmt.Print("'>")
			fmt.Println("        Your browser does not support the audio element.")
			fmt.Println("      </audio>")

			fmt.Println("</td>")
			fmt.Println("</tr>")

			if len(f.Title) != 0 {
				filenameTitle[f.Name] = f.Title
			}
		}
	}
}

func writeTopTitle(title, creator, year, id string, subjects []string) {
	//fmt.Println(year, "-", title, "--", creator, id)
	fmt.Printf("%s <a href=\"https://archive.org/details/%s\">%s</a> - %s - %d - %s --- %s\n", year, id, title, creator, len(subjects), subjects, id)
}

// There isn't always a title, so use 1) the name of the original (if exists); or 2) the filename;
// cache the filename
func makeFileTitle(title, name string, original []string, filenameTitle map[string]string) string {
	if len(title) != 0 {
		return title
	}

	if len(original) != 0 && len(original[0]) != 0 {
		if title, ok := filenameTitle[original[0]]; ok {
			return title
		}
	}

	return name

}

func findBaseMP3Filenames(files *[]ia.File) []string {
	var fns []string

	for i := 0; i < len(*files); i++ {
		f := (*files)[i]
		if f.Format == VBR_MP3 {
			name, _ := strings.CutSuffix(f.Name, VBR_MP3_SUFFIX)
			fns = append(fns, name)
			continue
		}
	}
	return fns
}

func collectCopies(files []ia.File) map[string][]*ia.File {
	nameCopies := make(map[string][]*ia.File)

	for i := 0; i < len(files); i++ {
		f := files[i]
		if _, ok := FileFormats[f.Format]; ok {
			// Assumes JSON ia.File ordering reflects track ordering. Might not be true for all
			f.TrackOrder = i
			baseName := makeFileBaseName(f.Name, f.Format)

			var copyFiles []*ia.File
			var ok bool
			if copyFiles, ok = nameCopies[baseName]; !ok {
				copyFiles = []*ia.File{&f}

			} else {
				copyFiles = append(copyFiles, &f)
			}
			nameCopies[baseName] = copyFiles
		}
	}

	return nameCopies
}

func findWantedCopies(copies map[string][]*ia.File, wantedFormats []string) []*ia.File {
	var wantedCopies []*ia.File

	for _, files := range copies {
		wantedCopy := findWantedCopy(files, wantedFormats)
		if wantedCopy != nil {
			wantedCopies = append(wantedCopies, wantedCopy)
		}
	}
	return wantedCopies
}

func findWantedCopy(files []*ia.File, wantedFormats []string) *ia.File {
	for _, format := range wantedFormats {
		for _, file := range files {
			if file.Format == format {
				return file
			}
		}
	}
	return nil
}

func makeFileBaseName(s string, format string) string {
	if ext, ok := FileFormats[format]; ok {
		return strings.TrimSuffix(s, ext)
	}

	return s

}

func comment(s string) {
	fmt.Printf("<!-- %s -->\n", s)
}
