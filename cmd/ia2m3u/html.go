package main

import (
	"cmp"
	"fmt"
	ia "github.com/gnewton/iascrape"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

func simpleHTML(item *ia.ItemTopLevelMetadata, wantedFormats []string, verbose bool) {

	copies := collectCopies(item.Files)

	wantedCopies := findWantedCopies(copies, wantedFormats)

	slices.SortFunc(wantedCopies, tunesByTrackOrder)

	meta := item.Metadata

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("<tr>")
	//fmt.Println("<td rowspan='", 11, "' valign='top' align='right'>")
	//fmt.Printf("<td rowspan='%d' valign='top' align='right'>\n", countAudioFiles(item.Files)+1)
	fmt.Printf("<td rowspan='%d' valign='top' align='right'>\n", len(wantedCopies)+1)
	//fmt.Println("<td valign='top' align='right'>")
	thumb := "https://" + item.D1 + item.Dir + "/" + Thumb

	if has, jp2f := hasJP2ZipFile(item.Files, item.Metadata.Identifier); !has {
		// LP FRONT cover image
		fmt.Println("<a href=\"https://" + item.D1 + item.Dir + "/" + item.Metadata.Identifier + "_itemimage.jpg\">")
		fmt.Println("<img width=160 align='left'   style='float: left;'    src='" + thumb + "'>")
		fmt.Println("</a>")
	} else {
		// LP FRONT cover image
		fmt.Println("<a href=\"" + makeJP2ImageUrl(jp2f, item, "0") + "\">")
		fmt.Println("<img width=160 align='left'   style='float: left;'    src='" + thumb + "'>")
		fmt.Println("</a>")

		// LP BACK cover image
		jp2ImageUrl := makeJP2ImageUrl(jp2f, item, "1")
		fmt.Print("<br><br>  &emsp; ")
		fmt.Println("<a href=\"" + jp2ImageUrl + "\">")
		fmt.Println("<img width=160 align='left'   style='float: left;'    src='" + jp2ImageUrl + "'>")
		fmt.Println("</a>")
	}
	fmt.Println("</td>")

	fmt.Println("")
	fmt.Println("<td valign='top' colspan='2'>")
	fmt.Println("<b>")

	var year string
	if meta.CanonicalYear == 0 {
		year = "[Unknown year]"
	} else {
		year = strconv.Itoa(meta.CanonicalYear)
	}

	title, creator := makeTitleCreator(meta.Titles, meta.Creators)
	//fmt.Printf("%s <a href=\"https://archive.org/details/%s\">%s</a> - %s - %d - %s --- %s\n", year, meta.Identifier, title, creator, len(meta.Subjects), meta.Subjects, meta.Identifier)

	fmt.Printf("%s <a href=\"https://archive.org/details/%s\">%s</a> - %s\n", year, meta.Identifier, title, creator)

	if verbose {
		fmt.Printf(" <a href=\"https://archive.org/metadata/%s\">JSON</a>\n", meta.Identifier)
	}

	fmt.Println("</td>")
	fmt.Println("</tr>")

	if len(item.Files) > 0 {
		writeAudioFiles(wantedCopies, meta.Identifier, verbose)
	}

	fmt.Println("<tr>  <td colspan='3'> <hr> </td> </tr>")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
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
		if _, ok := FileFormats[f.Format]; ok {
			n++
			fmt.Println("")
			fmt.Println("<tr valign='top'>")
			fmt.Println("<td width='35%'>")
			fmt.Printf("%d.\n", n)
			fmt.Printf("<a href=\"%s\">%s</a>", makeRemoteAudioURL(id, f.Name), makeFileTitle(f.Title, f.Name, f.Original, filenameTitle))
			//fmt.Printf("<a href=\"%s\">%s</a> %s", makeRemoteAudioURL(id, f.Name), makeFileTitle(f.Title, f.Name, f.Original, filenameTitle), f.Format)
			if verbose {
				fmt.Println(f.TrackOrder, f.Format)
			}

			fmt.Println("</td>")

			fmt.Println("<td>")
			// fmt.Println("<br>")
			fmt.Println("<p>")
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
