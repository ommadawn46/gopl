package github

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
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
