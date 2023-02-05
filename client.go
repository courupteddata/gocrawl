package gocrawl

import (
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient http.Client
	logger     *log.Logger
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{},
		logger:     log.Default(),
	}
}

func (c *Client) Crawl(url url.URL) []string {
	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		c.logger.Printf("Failed to create new request. err=%v", err)
		return nil
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.logger.Printf("Failed to make request. err=%v", err)
		return nil
	}

	c.extractLinks(response.Body, url)

	return []string{"empty"}
}

func (c *Client) extractLinks(reader io.Reader, baseUrl url.URL) {
	tokenizer := html.NewTokenizer(reader)

	urlMap := make(map[string]bool)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			// End of page
			return
		case html.StartTagToken:
			token := tokenizer.Token()
			for _, attribute := range token.Attr {
				if attribute.Key == "href" {
					href := attribute.Val
					parsedHref, err := url.Parse(href)
					if err != nil {
						c.logger.Printf("Failed to parse URL original: %v parsed: %v", href, parsedHref)
						continue
					}

					resolvedUrl := baseUrl.ResolveReference(parsedHref)
					c.logger.Printf("ResolvedUrl: %v, parsedHref: %v", resolvedUrl.String(), parsedHref)
					if resolvedUrl.Scheme == "https" || resolvedUrl.Scheme == "http" {
						urlMap[resolvedUrl.String()] = true
					}
				}
			}
		}
	}
}
