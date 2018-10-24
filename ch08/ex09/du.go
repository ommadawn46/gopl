package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Process struct {
	rootName       string
	fileSizes      chan int64
	nfiles, nbytes int64
	fin            bool
}

var vFlag = flag.Bool("v", false, "show verbose progress messages")

func main() {
	flag.Parse()

	rootDir := "/"
	roots, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	var processWait sync.WaitGroup
	finish := make(chan struct{})
	var processes []*Process
	for _, root := range roots {
		if !root.IsDir() {
			continue
		}
		processWait.Add(1)
		p := Process{rootDir + root.Name(), make(chan int64), 0, 0, false}
		processes = append(processes, &p)

		var n sync.WaitGroup
		n.Add(1)

		go walkDir(p.rootName, &n, p.fileSizes)
		go func(p *Process) {
			for size := range p.fileSizes {
				p.nfiles++
				p.nbytes += size
			}
			p.fin = true
		}(&p)
		go func() {
			n.Wait()
			close(p.fileSizes)
			processWait.Done()
		}()
	}

	go func() {
		processWait.Wait()
		finish <- struct{}{}
	}()

	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}
loop:
	for {
		select {
		case <-finish:
			break loop
		case <-tick:
			printDiskUsage(processes)
		}
	}
	printDiskUsage(processes)
}

func printDiskUsage(processes []*Process) {
	fmt.Println("------------------------------------------------")
	for _, p := range processes {
		check := " "
		if p.fin {
			check = "x"
		}
		fmt.Printf("[%s] %-16s %d files\t%.1f GB\n", check, p.rootName, p.nfiles, float64(p.nbytes)/1e9)
	}
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
