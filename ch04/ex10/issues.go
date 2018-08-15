package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	olderYearItems := []github.Issue{}
	olderMonthItems := []github.Issue{}
	newerMonthItems := []github.Issue{}

	for _, item := range result.Items {
		oneYearAgo := time.Now().AddDate(-1, 0, 0)
		oneMonthAgo := time.Now().AddDate(0, -1, 0)
		if oneYearAgo.After(item.CreatedAt) {
			olderYearItems = append(olderYearItems, *item)
		} else if oneMonthAgo.After(item.CreatedAt) {
			olderMonthItems = append(olderMonthItems, *item)
		} else {
			newerMonthItems = append(newerMonthItems, *item)
		}
	}

	fmt.Printf("%d issues (older than 1 year):\n", len(olderYearItems))
	for _, item := range olderYearItems {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}

	fmt.Printf("\n%d issues (older than 1 month):\n", len(olderMonthItems))
	for _, item := range olderMonthItems {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}

	fmt.Printf("\n%d issues (newer than 1 month):\n", len(newerMonthItems))
	for _, item := range newerMonthItems {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}
