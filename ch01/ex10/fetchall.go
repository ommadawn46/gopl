package main

import (
  "fmt"
  "io"
  "net/http"
  "os"
  "time"
  "strings"
)

func main() {
  start := time.Now()
  ch := make(chan string)
  for _, url := range os.Args[1:] {
    go fetch(url, ch)
  }
  for range os.Args[1:] {
    fmt.Println(<-ch)
  }
  fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string){
  start := time.Now()
  resp, err := http.Get(url)
  if err != nil {
    ch <- fmt.Sprint(err)
    return
  }
  defer resp.Body.Close()

  split_url := strings.Split(url, "/")
  file_name := split_url[len(split_url)-1]
  file, err := os.Create(file_name)
  if err != nil {
    ch <- fmt.Sprintf("while creating %s: %v", file_name, err)
    return
  }
  defer file.Close()

  nbytes, err := io.Copy(file, resp.Body)
  if err != nil {
    ch <- fmt.Sprintf("while reading %s: %v", url, err)
    return
  }
  secs := time.Since(start).Seconds()

  ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
