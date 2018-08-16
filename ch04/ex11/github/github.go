package github

import (
	"fmt"
	"time"
)

const BaseIssuesURL = "https://api.github.com/repos/%s/%s/issues"

func GetIssuesURL(owner string, repo string) string {
	return fmt.Sprintf(BaseIssuesURL, owner, repo)
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	REPOURL   string `json:"repository_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type IssuePost struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type IssuePatch struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	State string `json:"state"`
}
