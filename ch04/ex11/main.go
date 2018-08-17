package main

import (
	"./github"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func editWithEditor(inits string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "tmp")
	if err != nil {
		return inits, err
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	tmpfile.Write([]byte(inits))

	cmd := exec.Command(os.Getenv("EDITOR"), tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return inits, err
	}

	if out, err := ioutil.ReadFile(tmpfile.Name()); err != nil {
		return inits, err
	} else {
		return string(out), nil
	}
}

func printIssue(issue github.Issue) {
	repo := strings.TrimPrefix(issue.REPOURL, "https://api.github.com/repos/")
	fmt.Printf("\n--======================================================--\n"+
		"Repository: %s\nNumber: %d\nUser: %s\nTitle: %s\nCreatedAt: %s\nState: %s\n"+
		"----------------------------------------------------------\n\n"+
		"%s\n\n",
		repo, issue.Number, issue.User.Login, issue.Title, issue.CreatedAt, issue.State, issue.Body)
}

func read(owner string, repo string, number int) error {
	if number <= 0 {
		issues, err := github.GetIssues(owner, repo)
		if err != nil {
			return err
		}
		for _, issue := range issues {
			printIssue(issue)
		}
	} else {
		issue, err := github.GetIssue(owner, repo, number)
		if err != nil {
			return err
		}
		printIssue(issue)
	}
	return nil
}

func create(owner string, repo string) error {
	title, err := editWithEditor("input title here")
	if err != nil {
		return err
	}

	body, err := editWithEditor("input body here")
	if err != nil {
		return err
	}

	issue, err := github.CreateIssue(owner, repo, title, body)
	if err != nil {
		return err
	}

	fmt.Println("\nSuccessfully created Issue")
	printIssue(issue)
	return nil
}

func edit(owner string, repo string, number int) error {
	issue, err := github.GetIssue(owner, repo, number)
	if err != nil {
		return err
	}

	title, err := editWithEditor(issue.Title)
	if err != nil {
		return err
	}

	body, err := editWithEditor(issue.Body)
	if err != nil {
		return err
	}

	issue, err = github.EditIssue(owner, repo, number, title, body, true)
	if err != nil {
		return err
	}

	fmt.Println("\nSuccessfully edited Issue")
	printIssue(issue)
	return nil
}

func close(owner string, repo string, number int) error {
	issue, err := github.GetIssue(owner, repo, number)
	if err != nil {
		return err
	}

	issue, err = github.EditIssue(owner, repo, number, issue.Title, issue.Body, false)
	if err != nil {
		return err
	}

	fmt.Println("\nSuccessfully closed Issue")
	printIssue(issue)
	return nil
}

func main() {
	actionFlag := flag.String("a", "read", "Action to Issues (read, create, edit, close)")
	ownerFlag := flag.String("o", "", "Repository Owner (Required)")
	repoFlag := flag.String("r", "", "Repository (Required)")
	numberFlag := flag.Int("n", -1, "Issue number")
	flag.Parse()

	action, owner, repo, number := *actionFlag, *ownerFlag, *repoFlag, *numberFlag
	if owner == "" || repo == "" {
		flag.Usage()
	} else {
		var err error
		switch action {
		case "read":
			err = read(owner, repo, number)
		case "create":
			err = create(owner, repo)
		case "edit":
			err = edit(owner, repo, number)
		case "close":
			err = close(owner, repo, number)
		default:
			flag.Usage()
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
