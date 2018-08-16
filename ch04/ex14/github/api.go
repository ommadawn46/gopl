package github

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetIssue(owner string, repo string, number int) (int, Issue, error) {
	var issue Issue

	url := GetIssuesURL(owner, repo) + "/" + strconv.Itoa(number)
	resp, err := http.Get(url)
	if err != nil {
		return 0, issue, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return 0, issue, err
	}

	return resp.StatusCode, issue, nil
}

func GetIssues(owner string, repo string) (int, []Issue, error) {
	var issues []Issue

	url := GetIssuesURL(owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return 0, issues, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return 0, issues, err
	}

	return resp.StatusCode, issues, nil
}
