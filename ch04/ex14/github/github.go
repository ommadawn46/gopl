package github

import (
	"fmt"
	"time"
)

const BaseIssuesURL = "https://api.github.com/repos/%s/%s/issues?per_page=100"

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
	Milestone	*Milestone
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Milestone struct {
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	Desctiption	string
	CreatedAt time.Time `json:"created_at"`
	Creator		*User
}
