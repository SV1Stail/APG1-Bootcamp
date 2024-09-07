package spider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func CrawlWeb(ctx context.Context, urls <-chan string) <-chan *string {
	res := make(chan *string)
	workers := 8
	var wg sync.WaitGroup

	sem := make(chan interface{}, workers)

	go func() {
		defer close(res)
		for url := range urls {
			select {
			case <-ctx.Done():
				fmt.Println("Crawling stopped by context cancellation")
				return
			case sem <- struct{}{}:
				wg.Add(1)
				go func(url string) {
					defer func() { <-sem }()
					defer wg.Done()
					body, err := fetchPage(ctx, url)
					if err != nil {
						fmt.Printf("Error fetching page %s: %v\n", url, err)
						return
					}
					select {
					case <-ctx.Done():
						fmt.Println("Sending result stopped by context cancellation")
					case res <- body:
					}
				}(url)
			}
		}
		wg.Wait()
	}()

	return res
}

func fetchPage(ctx context.Context, url string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 | status code: %d", resp.StatusCode)

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := string(body)
	return &res, nil
}
