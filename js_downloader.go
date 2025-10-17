package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("github.com/unvalidor")
		fmt.Println("Usage: go run js_downloader.go js_endpoints.txt")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	os.MkdirAll("js_files", os.ModePerm)

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	sem := make(chan struct{}, 50) 

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(url string) {
			defer wg.Done()
			downloadFile(url)
			<-sem
		}(url)
	}

	wg.Wait()
}

func downloadFile(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Skipped:", url, resp.Status)
		return
	}

	fileName := path.Base(url)
	out, err := os.Create("js_files/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Downloaded:", url)
}
