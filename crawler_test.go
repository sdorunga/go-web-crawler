package main

import (
  "testing"
  "reflect"
  "net/http/httptest"
  "net/http"
  "net/url"
  "fmt"
)

type TestPrinter struct {
  lines []string
  expected []string
}

func (this *TestPrinter) Println(line string) {
  this.lines = append(this.lines, line)
}

func (this *TestPrinter) verify() bool {
  return reflect.DeepEqual(this.lines, this.expected)
}

func TestPrintsEmptyOutputWhenNoLinks(t *testing.T) {
  site := "http://site.com"
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintln(w,
  `<!DOCTYPE html>
  <html>
    <body>
    </body>
  </html>`)
  }))
  transport := &http.Transport{
    Proxy: func(req *http.Request) (*url.URL, error) {
      return url.Parse(server.URL)
    },
  }
  httpClient := http.Client{Transport: transport}
  client := &Client{httpClient}
  printer := TestPrinter{}
  crawler := Crawler{printer: &printer, url: site, client: client}
  crawler.Crawl()
  printer.expected = []string{
    "Crawled " + site + " and found the following.",
    "Links:",
    "------",
    "Static Assets:",
    "--------------",
  }
  crawler.Print()

  if printer.verify() != true {
    t.Error("Expected: ", printer.expected, ". Got ", printer.lines)
  }
}

func TestPrintsCorrectLinksAndStaticAssetsOnOnePage(t *testing.T) {
  site := "http://site.com"
  printer := TestPrinter{}
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    if r.RequestURI == site + "/" {
      fmt.Fprintln(w,
      `<!DOCTYPE html>
      <html>
      <body>
      <a href="https://site.com/about">About</a>
      <a href="/blog">Blog</a>
      <a href="https://externalsite.com/about">About External</a>
      <img src="https://cdn.static/img/horse.png"
      </body>
      </html>`)
    } else {
      fmt.Fprintln(w, "")
    }
  }))
  transport := &http.Transport{
    Proxy: func(req *http.Request) (*url.URL, error) {
      return url.Parse(server.URL)
    },
  }
  httpClient := http.Client{Transport: transport}
  client := &Client{httpClient}
  crawler := Crawler{printer: &printer, url: site, client: client}
  crawler.Crawl()
  printer.expected = []string{
    "Crawled " + site + " and found the following.",
    "Links:",
    "------",
    site + ":",
    "    https://site.com/about\n    http://site.com/blog",
    "Static Assets:",
    "--------------",
    site + ":",
    "    https://cdn.static/img/horse.png",
  }
  crawler.Print()

  if printer.verify() != true {
    t.Error("Expected: ", printer.expected, ". Got ", printer.lines)
  }
}
