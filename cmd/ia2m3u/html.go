package main

import (
	"fmt"
	ia "github.com/gnewton/iascrape"
	"net/url"
)

func simpleHTML2(count int64, item *ia.ItemTopLevelMetadata) {
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
	fmt.Println("<td valign='top' colspan=4>")
	fmt.Println("<b>")

	title, creator := makeTitleCreator(meta.Titles, meta.Creators)
	fmt.Printf("%s <a href=\"https://archive.org/details/%s\">%s</a> - %s - %d - %s --- %s\n", meta.CanonicalYear, meta.Identifier, title, creator, len(meta.Subjects), meta.Subjects, meta.Identifier)
	fmt.Println("</td>")
	fmt.Println("</tr>")

	// fmt.Println("1971      &nbsp;")
	// fmt.Println("<a href=\"https://archive.org/details/lp_champions-of-the-world_the-edinburgh-police-pipe-band\">")
	// fmt.Println("Champions Of The World</a>")
	// fmt.Println("<em>")
	// fmt.Println("The Edinburgh Police Pipe Band")
	// fmt.Println("</em>")
	// fmt.Println("</b>")

	//writeTopTitle(title, creator, meta.CanonicalYear, meta.Identifier, item.Metadata.Subjects)
	//fmt.Println("<ul>")
	if len(item.Files) > 0 {
		writeAudioFiles2(item.Files, meta.Identifier)
	}
	fmt.Println("<tr>")
	fmt.Println("<td>")
	fmt.Println("&nbsp;  ")
	fmt.Println("</td colspan='3'>")
	fmt.Println("</tr>")

}

func countAudioFiles(files []ia.File) int {
	n := 0
	for i := 0; i < len(files); i++ {
		f := files[i]
		if _, ok := FileFormats[f.Format]; ok {
			n++
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
			n++
			fmt.Println("")
			fmt.Println("<tr valign='top'>")
			fmt.Println("<td width='40%'>")
			fmt.Println(n)
			fmt.Printf(". <a href=\"%s\">%s</a> %s", makeRemoteAudioURL(id, f.Name), makeFileTitle(f.Title, f.Name, f.Original, filenameTitle), f.Format)
			fmt.Println("</td>")

			fmt.Println("<td>")
			// fmt.Println("<br>")
			// fmt.Println("<br>")
			fmt.Println("      <audio controls>")
			fmt.Print("        <source preload='none' src=\"")
			fmt.Print(AudioFileBaseUrl + url.PathEscape(id) + "/" + url.PathEscape(f.Name))
			fmt.Print("\"")
			fmt.Print("  type='audio/mpeg'>")
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

func simpleHTML(count int64, item *ia.ItemTopLevelMetadata) {

	if true {
		simpleHTML2(count, item)
		return
	}
	fmt.Println("<li>")

	thumb := "https://" + item.D1 + item.Dir + "/" + Thumb

	if has, jp2f := hasJP2ZipFile(item.Files, item.Metadata.Identifier); !has {
		// LP FRONT cover image
		fmt.Println("	    <a href=\"")
		fmt.Println("https://" + item.D1 + item.Dir + "/" + item.Metadata.Identifier + "_itemimage.jpg\">")

		fmt.Println("      <img width=160 align='left'   style='float: left;'    src='" + thumb + "'>")
		fmt.Println("	    </a>")
	} else {
		// LP FRONT cover image
		fmt.Println("	    <a href=\"" + makeJP2ImageUrl(jp2f, item, "0") + "\">")
		fmt.Println("      <img width=160 align='left'   style='float: left;'    src='" + thumb + "'>")
		fmt.Println("	    </a>")

		// LP BACK cover image
		jp2ImageUrl := makeJP2ImageUrl(jp2f, item, "1")
		fmt.Print("	    <br><br>  &emsp; ")
		fmt.Println("<a href=\"" + jp2ImageUrl + "\">")
		fmt.Println("      <img width=160 align='left'   style='float: left;'    src='" + jp2ImageUrl + "'>")
		fmt.Println("	    </a>")
	}

	meta := item.Metadata
	title, creator := makeTitleCreator(meta.Titles, meta.Creators)
	writeTopTitle(title, creator, meta.CanonicalYear, meta.Identifier, item.Metadata.Subjects)
	fmt.Println("<ul>")
	if len(item.Files) > 0 {
		writeAudioFiles(item.Files, meta.Identifier)
	}
	fmt.Println("</ul>")
	fmt.Println("</li>")
}

func writeTopTitle(title, creator, year, id string, subjects []string) {
	//fmt.Println(year, "-", title, "--", creator, id)
	fmt.Printf("%s <a href=\"https://archive.org/details/%s\">%s</a> - %s - %d - %s --- %s\n", year, id, title, creator, len(subjects), subjects, id)
}

func writeAudioFiles(files []ia.File, id string) {
	filenameTitle := make(map[string]string)
	for i := 0; i < len(files); i++ {
		f := files[i]
		if _, ok := FileFormats[f.Format]; ok {

			fmt.Println("<li>")

			fmt.Printf("<a href=\"%s\">%s</a> %s", makeRemoteAudioURL(id, f.Name), makeFileTitle(f.Title, f.Name, f.Original, filenameTitle), f.Format)

			fmt.Println("      <audio controls>")
			fmt.Print("        <source preload='none' src=\"")
			fmt.Print(AudioFileBaseUrl + url.PathEscape(id) + "/" + url.PathEscape(f.Name))
			fmt.Print("\"")
			fmt.Print("  type='audio/mpeg'>")
			fmt.Println("        Your browser does not support the audio element.")
			fmt.Println("      </audio>")

			fmt.Println("</li>")

			if len(f.Title) != 0 {
				filenameTitle[f.Name] = f.Title
			}
		}
	}
}

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
