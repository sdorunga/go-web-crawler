package main

import (
  "testing"
  "reflect"
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
  url := "http://crawler.net"
  printer := TestPrinter{}
  crawler := Crawler{&printer, url}
  crawler.Crawl(url)
  printer.expected = []string{
    "Crawled " + url + " and found the following.",
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
