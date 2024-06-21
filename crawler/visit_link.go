package crawler

import (
	"fmt"
	"gowebcrawler/db"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var (
	Reset = "\033[0m"
	Green = "\033[32m"
	Blue  = "\033[34m"
	Red   = "\033[31m"
)

type VisitedLink struct {
	Website     string    `bson:"website"`
	Link        string    `bson:"link"`
	VisitedDate time.Time `bson:"visited_data"`
}

type ErrorLog struct {
	Link        string    `bson:"link"`
	VisitedDate time.Time `bson:"visited_data"`
	Error       string    `bson:"error"`
}

var (
	visitedLinks = make(map[string]bool)
	mu           sync.Mutex
)

func Start(initialLink string) {
	linkChan := make(chan string)
	doneChan := make(chan bool)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		linkChan <- initialLink
	}()

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	go func() {
		for link := range linkChan {
			wg.Add(1)
			go func(link string) {
				defer wg.Done()
				visitLink(link, linkChan)
			}(link)
		}
	}()

	<-doneChan
}

func visitLink(link string, linkChan chan string) {
	fmt.Printf("Visitando: %s\n", link)

	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("[error] Error on http.Get: %s\n", err)
		logError(link, fmt.Sprintf("Error on http.Get: %s", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s[error] Status diferente de 200: %d%s\n", Red, resp.StatusCode, Reset)
		logError(link, fmt.Sprintf("Status diferente de 200: %d", resp.StatusCode))
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("%s[error] Error on html.Parse: %s%s\n", Red, err, Reset)
		logError(link, fmt.Sprintf("Error on html.Parse: %s", err))
		return
	}

	extractLinks(doc, linkChan)
}

func extractLinks(node *html.Node, linkChan chan string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			link, err := url.Parse(attr.Val)
			if err != nil || link.Scheme == "" || link.Scheme == "mailto" {
				continue
			}

			if link.Scheme != "http" && link.Scheme != "https" {
				fmt.Println("Ignorando URL com esquema não suportado:", link)
				continue
			}

			if link.Fragment != "" {
				fmt.Println("Ignorando URL com âncora:", link)
				continue
			}

			mu.Lock()
			if db.VisitedLink(link.String()) {
				mu.Unlock()
				fmt.Printf("%sLink já visitado %s%s\n", Blue, link, Reset)
				continue
			}
			mu.Unlock()

			visitedLink := VisitedLink{
				Website:     link.Host,
				Link:        link.String(),
				VisitedDate: time.Now(),
			}

			mu.Lock()
			db.Insert("links", visitedLink)
			mu.Unlock()

			fmt.Printf("%sLink visitado %s%s\n", Green, link, Reset)
			linkChan <- link.String()
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, linkChan)
	}
}

func logError(link, message string) {
	mu.Lock()
	defer mu.Unlock()
	errorLog := ErrorLog{
		Link:        link,
		VisitedDate: time.Now(),
		Error:       message,
	}
	db.Insert("error_log", errorLog)
}
