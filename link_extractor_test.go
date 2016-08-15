
package main

import (
  "testing"
  "reflect"
)

type TestClient struct {
  response string
}

func (this *TestClient) Get(url string) string {
  return this.response
}

const (
  twoUrls = `
  <!DOCTYPE html>
  <html>
    <body>
      <a href="https://site.com/about">About</a>
      <a href="/blog">Blog</a>
      <a href="https://externalsite.com/about">About External</a>
    </body>
  </html> `
  assetUrls = `
  <!DOCTYPE html>
  <html>
    <head>
      <script src="myscript.js"></script>
      <link rel="stylesheet" type="text/css" href="mystyle.css">
    </head>
    <body>
      <img src="https://site.com/img/horse.png"/>
    </body>
  </html>`
)

func TestReturnsAListOfLinksFromHttp(t *testing.T) {
  site := "https://site.com"
  response := assetUrls
  extractor := LinkExtractor{&TestClient{response}, site}
  links := extractor.getLinks().ResourceLinks
  expectedLinks := []string{"https://site.com/myscript.js","https://site.com/mystyle.css", "https://site.com/img/horse.png"}
  if !reflect.DeepEqual(links, expectedLinks) {
    t.Error("Expected: ", expectedLinks, ". Got: ", links)
  }
}
