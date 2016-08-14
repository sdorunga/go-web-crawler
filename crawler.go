package main

import (
  "os"
  "fmt"
)

func main() {
  url := os.Args[1]
  crawler := Crawler{&TerminalPrinter{}, url}
  crawler.Crawl(url)
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
}

func (this *Crawler) Crawl(url string) {
}

func (this *Crawler) Print() {
  this.printer.Println("Crawled " + this.url + " and found the following.")
  this.printer.Println("Links:")
  this.printer.Println("------")
  this.printer.Println("Static Assets:")
  this.printer.Println("--------------")
}


