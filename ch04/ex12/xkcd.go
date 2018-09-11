package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Info struct {
	Num        int
	Year       string
	Month      string
	Day        string
	Img        string
	Title      string
	Link       string
	News       string
	Safe_title string
	Transcript string
	Alt        string
}

const DBPath = "./xkcd.dat"

func serialize(infos []Info) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(&infos)
	return buf.Bytes(), err
}

func deserialize(data []byte) ([]Info, error) {
	var infos []Info
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&infos)
	return infos, err
}

func saveInfos(infos []Info) error {
	if data, err := serialize(infos); err != nil {
		return err
	} else {
		return ioutil.WriteFile(DBPath, data, 0644)
	}
}

func loadInfos() ([]Info, error) {
	if data, err := ioutil.ReadFile(DBPath); err != nil {
		return []Info{}, err
	} else {
		return deserialize(data)
	}
}

func getInfoByNum(num int) (Info, error) {
	numStr := ""
	if num > 0 {
		numStr = strconv.Itoa(num)
	}
	url := fmt.Sprintf("https://xkcd.com/%s/info.0.json", numStr)

	resp, err := http.Get(url)
	if err != nil {
		return Info{}, err
	}
	defer resp.Body.Close()

	var info Info
	return info, json.NewDecoder(resp.Body).Decode(&info)
}

func getRemoteHeadNum() int {
	if info, err := getInfoByNum(0); err == nil {
		return info.Num
	} else {
		return 0
	}
}

func getSavedHeadNum(infos []Info) int {
	headNum := 0
	for _, info := range infos {
		if info.Num > headNum {
			headNum = info.Num
		}
	}
	return headNum
}

func fileIsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func downloadInfos(startNum int, endNum int) []Info {
	var infos []Info
	for num := startNum; num <= endNum; num++ {
		indicator := []string{"/", "-", "\\", "|"}
		fmt.Printf("[%s] Downloading comic information .. %d / %d\r", indicator[num%4], (num-startNum)+1, (endNum-startNum)+1)
		info, _ := getInfoByNum(num)
		infos = append(infos, info)
	}
	fmt.Println("")
	return infos
}

func searchInfos(comicInfos []Info, query []string) []Info {
	var hits []Info
	containsAll := func(content string) bool {
		for _, q := range query {
			if !strings.Contains(content, q) {
				return false
			}
		}
		return true
	}
	for _, info := range comicInfos {
		content := strings.Join([]string{strconv.Itoa(info.Num), info.Year, info.Month,
			info.Day, info.Img, info.Title, info.Link, info.News, info.Safe_title, info.Transcript, info.Alt}, "/")
		if containsAll(content) {
			hits = append(hits, info)
		}
	}
	return hits
}

func printInfos(infos []Info) {
	for _, info := range infos {
		fmt.Printf("\n--======================================================--\n"+
			"URL: %s\nTitle: %s\nNum: %d\nPost Date: %s/%s/%s\n"+
			"----------------------------------------------------------\n\n"+
			"%s\n\n",
			info.Img, info.Title, info.Num, info.Year, info.Month, info.Day, info.Transcript)
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s word1 word2 word3..\n", os.Args[0])
		return
	}

	var comicInfos []Info
	if fileIsExist(DBPath) {
		comicInfos, _ = loadInfos()
	}

	savedHeadNum := getSavedHeadNum(comicInfos)
	remoteHeadNum := getRemoteHeadNum()
	if savedHeadNum < remoteHeadNum {
		newInfos := downloadInfos(savedHeadNum+1, remoteHeadNum)
		comicInfos = append(comicInfos, newInfos...)
		saveInfos(comicInfos)
	}

	hits := searchInfos(comicInfos, os.Args[1:])
	printInfos(hits)
}
