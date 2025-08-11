package fetch

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Result летит из воркеров в главную горутину.
type Result struct {
	URL    string
	Status int
	Length int
	Links  int
	Err    error
}

// Worker скачивает страницу, выдёргивает ссылки, кладёт их обратно в jobs.
func Worker(id int, jobs chan string, results chan<- Result, wg *sync.WaitGroup, done <-chan struct{}, visited *Visited) {
	defer wg.Done()

	client := &http.Client{Timeout: 5 * time.Second}

	for u := range jobs {
		if visited.Seen(u) { // ➋ пропускаем дубликаты
			continue
		}
		body, st, err := fetch(client, u)
		if err != nil {
			results <- Result{URL: u, Status: st, Err: err}
			continue
		}

		found := extractLinks(body, u, jobs, done, visited)
		results <- Result{URL: u, Status: st, Length: len(body), Links: found}
	}
}

// fetch запрашивает страницу и возвращает тело.
func fetch(c *http.Client, raw string) ([]byte, int, error) {
	resp, err := c.Get(raw)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return data, resp.StatusCode, nil
}

// extractLinks находит <a href>, кидает внутрь домена в jobs, возвращает количество новых URL.
func extractLinks(body []byte, base string, jobs chan<- string, done <-chan struct{}, visited *Visited) int {
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return 0
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return 0
	}
	var found int
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					href := strings.TrimSpace(a.Val)
					link, err := baseURL.Parse(href)
					if err == nil && sameDomain(baseURL, link) && !visited.Seen(link.String()) {
						select {
						case jobs <- link.String():
							found++
						case <-done: // получили сигнал остановиться
							return
						default: // jobs полон — backpressure
						}
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return found
}

// sameDomain true, если host совпадает.
func sameDomain(a, b *url.URL) bool {
	return strings.EqualFold(a.Hostname(), b.Hostname())
}
