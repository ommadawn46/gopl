package main

import (
  "html/template"
  "net/http"
  "log"
  "./github"
)

var issueListTemplate = template.Must(template.New("issuelist").Parse(`
<h1>{{.Items | len}} issues</h1>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>User</th>
  <th>Title</th>
  <th>Milestone</th>
</tr>
{{range .Items}}
<tr>
  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
  {{ if .Milestone }}
    <td><a href='{{.Milestone.HTMLURL}}'>{{.Milestone.Title}}</a></td>
  {{ else }}
    <td>None</td>
  {{ end }}
</tr>
{{end}}
</table>
`))

type IssuesResult struct {
	Items         []github.Issue
}

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    var owner, repo string
    if err := r.ParseForm(); err != nil {
      log.Print(err)
    } else {
      if params, ok := r.Form["owner"]; ok {
        owner = params[0]
      }
      if params, ok := r.Form["repo"]; ok {
        repo = params[0]
      }
    }
    issues, _ := github.GetIssues(owner, repo)

    issuesResult := IssuesResult{issues}
    if err := issueListTemplate.Execute(w, issuesResult); err != nil {
      log.Print(err)
    }
  })
  http.ListenAndServe("localhost:8000", nil)
}
