package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var cancelCh = make(chan struct{}, 1)

func main() {
	var n sync.WaitGroup
	responses := make(chan *[]byte, 3)

	for _, url := range os.Args[1:] {
		n.Add(1)
		go func(url string) {
			defer n.Done()
			res, err := request(url)
			if err != nil {
				log.Printf(err.Error())
				return
			}
			responses <- res
			go func() {
				for {
					cancelCh <- struct{}{}
				}
			}()
		}(url)
	}
	n.Wait()

	close(responses)
	res := <-responses
	if res != nil {
		fmt.Printf("%s", res)
	}
}

func request(url string) (*[]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Cancel = cancelCh
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: %v\n", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("fetch: reading %s: %v\n", url, err)
	}
	return &b, nil
}
