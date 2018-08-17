package github

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"fmt"
	"strconv"
	"strings"
)

func GetIssue(owner string, repo string, number int) (Issue, error) {
	var issue Issue

	url := GetIssuesURL(owner, repo) + "/" + strconv.Itoa(number)
	resp, err := http.Get(url)
	if err != nil {
		return issue, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return issue, fmt.Errorf("%d: Failed to get Issue\n", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return issue, err
	}

	return issue, nil
}

func GetIssues(owner string, repo string) ([]Issue, error) {
	var issues []Issue

	url := GetIssuesURL(owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return issues, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return issues, fmt.Errorf("%d: Failed to get Issues\n", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return issues, err
	}

	return issues, nil
}

func CreateIssue(owner string, repo string, title string, body string) (Issue, error) {
	var issue Issue

	issuePost := IssuePost{strings.TrimSpace(title), strings.TrimSpace(body)}
	jsonBytes, err := json.Marshal(issuePost)
	if err != nil {
		return issue, err
	}

	client := new(http.Client)
	url := GetIssuesURL(owner, repo)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	req.SetBasicAuth(os.Getenv("GH_USER"), os.Getenv("GH_PASS"))
	resp, err := client.Do(req)
	if err != nil {
		return issue, err
	}

	if resp.StatusCode != 201 {
		return issue, fmt.Errorf("%d: Failed to create Issue\n", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return issue, err
	}

	return issue, nil
}

func EditIssue(owner string, repo string, number int,
	title string, body string, open bool) (Issue, error) {
	var issue Issue

	state := "open"
	if !open {
		state = "closed"
	}

	issuePatch := IssuePatch{strings.TrimSpace(title), strings.TrimSpace(body), state}
	jsonBytes, err := json.Marshal(issuePatch)
	if err != nil {
		return issue, err
	}

	client := new(http.Client)
	url := GetIssuesURL(owner, repo) + "/" + strconv.Itoa(number)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(jsonBytes))
	req.SetBasicAuth(os.Getenv("GH_USER"), os.Getenv("GH_PASS"))
	resp, err := client.Do(req)
	if err != nil {
		return issue, err
	}

	if resp.StatusCode != 200 {
		return issue, fmt.Errorf("%d: Failed to edit Issue\n", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return issue, err
	}
	return issue, nil
}
