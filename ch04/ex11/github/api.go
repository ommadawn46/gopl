package github

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func CreateIssue(owner string, repo string, title string, body string) (int, Issue, error) {
	var issue Issue

	issuePost := IssuePost{strings.TrimSpace(title), strings.TrimSpace(body)}
	jsonBytes, err := json.Marshal(issuePost)
	if err != nil {
		return 0, issue, err
	}

	client := new(http.Client)
	url := GetIssuesURL(owner, repo)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	req.SetBasicAuth(os.Getenv("GH_USER"), os.Getenv("GH_PASS"))
	resp, err := client.Do(req)
	if err != nil {
		return 0, issue, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return resp.StatusCode, issue, err
	}
	return resp.StatusCode, issue, nil
}

func EditIssue(owner string, repo string, number int,
	title string, body string, open bool) (int, Issue, error) {
	var issue Issue

	state := "open"
	if !open {
		state = "closed"
	}

	issuePatch := IssuePatch{strings.TrimSpace(title), strings.TrimSpace(body), state}
	jsonBytes, err := json.Marshal(issuePatch)
	if err != nil {
		return 0, issue, err
	}

	client := new(http.Client)
	url := GetIssuesURL(owner, repo) + "/" + strconv.Itoa(number)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(jsonBytes))
	req.SetBasicAuth(os.Getenv("GH_USER"), os.Getenv("GH_PASS"))
	resp, err := client.Do(req)
	if err != nil {
		return 0, issue, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return resp.StatusCode, issue, err
	}
	return resp.StatusCode, issue, nil
}
