package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"

	"../ex08/sorting"
)

var trackTable = template.Must(template.New("trackTable").Parse(`
<html>
   <body>
      <h1>Track Table</h1>
      <table border="1">
         <thead>
            <tr>
               <th><a href="/?sort=Title">Title</a></th>
               <th><a href="/?sort=Artist">Artist</a></th>
               <th><a href="/?sort=Album">Album</a></th>
               <th><a href="/?sort=Year">Year</a></th>
               <th><a href="/?sort=Length">Length</a></th>
            </tr>
         </thead>
         <tbody>
            {{range .}}
            <tr>
               <td>{{.Title}}</td>
               <td>{{.Artist}}</td>
               <td>{{.Album}}</td>
               <td>{{.Year}}</td>
               <td>{{.Length}}</td>
            </tr>
            {{end}}
         </tbody>
      </table>
   </body>
</html>
`))

func printTracks(w io.Writer, tracks []*sorting.Track) {
	if err := trackTable.Execute(w, tracks); err != nil {
		log.Fatal(err)
	}
}

func main() {
	tracks := []*sorting.Track{
		{"Go", "Delilah", "From the Roots Up", 2012, sorting.Length("3m38s")},
		{"Go", "Moby", "Moby", 1992, sorting.Length("3m37s")},
		{"Go Ahead", "Alicia Keys", "As I Am", 2007, sorting.Length("4m36s")},
		{"Ready 2 Go", "Martin Solveig", "Smash", 2011, sorting.Length("4m24s")},
	}
	multiSort := &sorting.MultiSort{Tracks: tracks}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sortby := r.FormValue("sort")
		switch sortby {
		case "Title":
			multiSort.AppendKey(sorting.TitleCompare)
		case "Year":
			multiSort.AppendKey(sorting.YearCompare)
		case "Length":
			multiSort.AppendKey(sorting.LengthCompare)
		case "Artist":
			multiSort.AppendKey(sorting.ArtistCompare)
		case "Album":
			multiSort.AppendKey(sorting.AlbumCompare)
		}
		sort.Sort(multiSort)
		printTracks(w, tracks)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
