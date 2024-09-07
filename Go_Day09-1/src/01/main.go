package main

import (
	"02/spider"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	urlChan := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Received interrupt signal, cancelling...")
		cancel()
	}()
	go func() {
		urls := []string{
			"https://www.youtube.com/",
		}

		for _, url := range urls {
			urlChan <- url
		}
		close(urlChan) //
	}()

	resChan := spider.CrawlWeb(ctx, urlChan)
	for r := range resChan {
		if r != nil {
			fmt.Println("Fetched page content:")
			fmt.Println(*r)
		}
	}
	fmt.Println("Crawling finished")
}
