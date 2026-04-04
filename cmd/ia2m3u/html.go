package main

import (
	"fmt"
	ia "github.com/gnewton/iascrape"
	"log"
	"net/url"
	"strconv"
)

func simpleHTML(item *ia.ItemTopLevelMetadata) {
	meta := item.Metadata
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("<tr>")
	//fmt.Println("<td rowspan='", 11, "' valign='top' align='right'>")
	fmt.Printf("<td rowspan='%d' valign='top' align='right'>\n", countAudioFiles(item.Files)+1)
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
	fmt.Printf("%s <a href=\"https://archive.org/details/%s\">%s</a> - %s - %d - %s --- %s\n", year, meta.Identifier, title, creator, len(meta.Subjects), meta.Subjects, meta.Identifier)

	fmt.Printf(" <a href=\"https://archive.org/metadata/%s\">JSON</a>\n", meta.Identifier)

	fmt.Println("</td>")
	fmt.Println("</tr>")

	if len(item.Files) > 0 {
		log.Println("%%%%%%%%%%%%%%%%%%%%%%% ", meta.Identifier)
		writeAudioFiles2(item.Files, meta.Identifier)
	}

	fmt.Println("<tr>  <td colspan='3'> <hr> </td> </tr>")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}

func countAudioFiles(files []ia.File) int {
	n := 0
	for i := 0; i < len(files); i++ {
		f := files[i]

		if _, ok := FileFormats[f.Format]; ok {
			n++
			//log.Println("+++++++++++ 567", f.Format)
		}
	}
	return n
}

func writeAudioFiles2(files []ia.File, id string) {

	filenameTitle := make(map[string]string)
	n := 0
	for i := 0; i < len(files); i++ {
		f := files[i]
		if _, ok := FileFormats[f.Format]; ok {
			log.Println(f.Format, f.Name, f.Title)
			n++
			fmt.Println("")
			fmt.Println("<tr valign='top'>")
			fmt.Println("<td width='35%'>")
			fmt.Printf("%d.\n", n)
			fmt.Printf("<a href=\"%s\">%s</a> %s", makeRemoteAudioURL(id, f.Name), makeFileTitle(f.Title, f.Name, f.Original, filenameTitle), f.Format)
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
			fmt.Println("<br>")
			fmt.Println("<br>")
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
