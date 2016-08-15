package main

import (
  "os"
  "fmt"
  "net/http"
  "io/ioutil"
  "strings"
)

func main() {
  url := os.Args[1]
  Client := &Client{http.Client{}}
  crawler := Crawler{&TerminalPrinter{}, url, Client, SiteLinks{}}
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

type HttpClient interface {
  Get(string) string
}

type Client struct {
  client http.Client
}

func (this *Client) Get(url string) string{
  resp, err := this.client.Get(url)
  if err != nil {
    fmt.Println("Error parsing url:", url)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Could not read body for request:", url)
  }
  return string(body)
}

type Crawler struct {
  printer Printer
  url string
  client HttpClient
  links SiteLinks
}

func (this *Crawler) Crawl() {
  extractor := LinkExtractor{this.client, this.url}
  this.links = extractor.getLinks()
}

func (this *Crawler) Print() {
  this.printer.Println("Crawled " + this.url + " and found the following.")
  this.printer.Println("Links:")
  this.printer.Println("------")
  if (len(this.links.PageLinks) > 0) {
    this.printer.Println(strings.Join(this.links.PageLinks, "\n"))
  }
  this.printer.Println("Static Assets:")
  this.printer.Println("--------------")
  if (len(this.links.PageLinks) > 0) {
    this.printer.Println(strings.Join(this.links.ResourceLinks, "\n"))
  }
}

