package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://martinfowler.com", "URL to build a sitemap for")
	maxDepth := flag.Int("depth", 1, "Depth of the sitemap")
	flag.Parse()

	links := bfs(*urlFlag, *maxDepth)

	toXml := urlset{
		Xmlns: xmlns,
	}
	for _, link := range links {
		toXml.Urls = append(toXml.Urls, loc{link})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("length", len(links))
}

func get(urlStr string) []string {
	fmt.Println("Getting", urlStr)
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // memory leak prevention

	baseUrl := &url.URL{
		Scheme: resp.Request.URL.Scheme,
		Host:   resp.Request.URL.Host,
	}
	links, _ := Parse(resp.Body)

	filter := func(link Link) bool {
		if strings.HasPrefix(link.Href, "/") {
			link.Href = urlStr + link.Href
		}

		url := strings.Split(link.Href, "?")[0]
		onlyHtml := strings.HasSuffix(url, ".html")
		sameDomain := strings.HasPrefix(link.Href, baseUrl.String())

		return onlyHtml && sameDomain
	}

	links = filterLinks(links, filter)
	links = removeDuplicates(links, baseUrl.String())

	var ret []string
	for _, link := range links {
		ret = append(ret, link.Href)
	}
	return ret
}

func filterLinks(links []Link, filter func(Link) bool) []Link {
	var ret []Link
	for _, link := range links {
		if filter(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func removeDuplicates(links []Link, baseUrl string) []Link {
	var ret []Link
	for _, link := range links {
		found := false
		link.Href = strings.Split(link.Href, "?")[0]
		if strings.HasPrefix(link.Href, "/") {
			link.Href = baseUrl + link.Href
		}
		for _, existing := range ret {
			if existing.Href == link.Href {
				found = true
				break
			}
		}
		if !found {
			ret = append(ret, link)
		}
	}
	return ret
}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				if _, ok := seen[link]; !ok {
					nq[link] = struct{}{}
				}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}
