package main

import (
  "flag"
  "fmt"
  "log"
  "io/ioutil"
  "os"
  "os/exec"
  "strings"
  "./github"
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
  if err := cmd.Run(); err != nil{
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
  fmt.Printf("\n---------------\nRepository: %s\nNumber: %d\nUser: %s\nTitle: %s\nCreatedAt: %s\nState: %s\n\n%s\n",
    repo, issue.Number, issue.User.Login, issue.Title, issue.CreatedAt, issue.State, issue.Body)
}

func read(owner string, repo string, number int) {
  if number < 0 {
    statusCode, issues, err := github.GetIssues(owner, repo)
    if err != nil {
      log.Fatal(err)
      return
    }
    if statusCode == 200 {
      for _, issue := range issues {
        printIssue(issue)
      }
    } else {
      fmt.Printf("%d: Failed to get Issue\n", statusCode)
    }
  } else {
    statusCode, issue, err := github.GetIssue(owner, repo, number)
    if err != nil {
      log.Fatal(err)
      return
    }
    if statusCode == 200 {
      printIssue(issue)
    } else {
      fmt.Printf("%d: Failed to get Issue\n", statusCode)
    }
  }
}

func create(owner string, repo string) {
  title, err := editWithEditor("input title here")
  if err != nil {
    log.Fatal(err)
    return
  }

  body, err := editWithEditor("input body here")
  if err != nil {
    log.Fatal(err)
    return
  }

  statusCode, issue, err := github.CreateIssue(owner, repo, title, body)
  if err != nil {
    log.Fatal(err)
    return
  }

  if statusCode == 201 {
    fmt.Printf("%d: Successfully created Issue\n\n", statusCode)
    printIssue(issue)
  } else {
    fmt.Printf("%d: Failed to create Issue\n", statusCode)
  }
}

func edit(owner string, repo string, number int) {
  statusCode, issue, err := github.GetIssue(owner, repo, number)
  if err != nil {
    log.Fatal(err)
    return
  }
  if statusCode != 200 {
    fmt.Printf("%d: Failed to get Issue\n", statusCode)
  }

  title, err := editWithEditor(issue.Title)
  if err != nil {
    log.Fatal(err)
    return
  }

  body, err := editWithEditor(issue.Body)
  if err != nil {
    log.Fatal(err)
    return
  }

  statusCode, issue, err = github.EditIssue(owner, repo, number, title, body, true)
  if err != nil {
    log.Fatal(err)
    return
  }

  if statusCode == 200 {
    fmt.Printf("%d: Successfully edited Issue\n\n", statusCode)
    printIssue(issue)
  } else {
    fmt.Printf("%d: Failed to edit Issue\n", statusCode)
  }
}

func close(owner string, repo string, number int) {
  statusCode, issue, err := github.GetIssue(owner, repo, number)
  if err != nil {
    log.Fatal(err)
    return
  }
  if statusCode != 200 {
    fmt.Printf("%d: Failed to get Issue\n", statusCode)
  }

  statusCode, issue, err = github.EditIssue(owner, repo, number, issue.Title, issue.Body, false)
  if err != nil {
    log.Fatal(err)
    return
  }

  if statusCode == 200 {
    fmt.Printf("%d: Successfully closed Issue\n\n", statusCode)
    printIssue(issue)
  } else {
    fmt.Printf("%d: Failed to close Issue\n", statusCode)
  }
}

func main() {
  actionFlag := flag.String("a", "read", "Action to Issues (read, create, edit, close)")
  ownerFlag := flag.String("o", "", "Repository Owner")
  repoFlag := flag.String("r", "", "Repository")
  numberFlag := flag.Int("n", -1, "Issue number")
  flag.Parse()

  action, owner, repo, number := *actionFlag, *ownerFlag, *repoFlag, *numberFlag
  if owner == "" || repo == "" {
    flag.Usage()
  } else {
    switch action {
    case "read": read(owner, repo, number)
    case "create": create(owner, repo)
    case "edit": edit(owner, repo, number)
    case "close": close(owner, repo, number)
    default: flag.Usage()
    }
  }
}
