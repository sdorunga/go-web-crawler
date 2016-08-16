package main

import (
  "os"
  "fmt"
  "strings"
)

func main() {
  url := os.Args[1]
  Client := NewClient()
  crawler := Crawler{printer: &TerminalPrinter{}, url: url, client: Client}
  crawler.Crawl()
  crawler.Print()
}

type Printer interface {
  Println(string)
}

type TerminalPrinter struct {
}

func (this *TerminalPrinter) Println(line string) {
  fmt.Println(line)
}

type Crawler struct {
  printer Printer
  url string
  client HttpClient
  links []SiteLinks
  seenLinks []string
}

func (this *Crawler) Crawl() {
  extractor := LinkExtractor{this.client, this.url}
  this.links = append(this.links, extractor.getLinks())
  this.seenLinks = append(this.seenLinks, this.url)

  this.crawlRelatedPages(extractor.getLinks().PageLinks)
}

func (this *Crawler) crawlRelatedPages(pageLinks []string) {
  if len(pageLinks) == 0 {
    return
  }
  for _, url := range pageLinks {
    if contains(this.seenLinks, url) {
      continue
    }
    this.seenLinks = append(this.seenLinks, url)
    fmt.Println(url)
    extractor := LinkExtractor{this.client, url}
    this.links = append(this.links, extractor.getLinks())
    this.crawlRelatedPages(extractor.getLinks().PageLinks)
  }
}

func (this *Crawler) Print() {
  this.printer.Println("Crawled " + this.url + " and found the following.")
  this.printer.Println("Links:")
  this.printer.Println("------")
  for _, link := range this.links {
    if (len(link.PageLinks) > 0) {
      this.printer.Println(link.Url + ":")
      this.printer.Println("    " + strings.Join(link.PageLinks, "\n    "))
    }
  }
  this.printer.Println("Static Assets:")
  this.printer.Println("--------------")
  for _, link := range this.links {
    if (len(link.ResourceLinks) > 0) {
      this.printer.Println(link.Url + ":")
      this.printer.Println("    " + strings.Join(link.ResourceLinks, "\n    "))
    }
  }
}

