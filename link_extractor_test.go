
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
  duplicates = `
  <!DOCTYPE html>
  <html>
    <head>
      <script src="myscript.js"></script>
      <script src="myscript.js"></script>
    </head>
    <body>
      <img src="https://site.com/img/horse.png"/>
      <img src="https://site.com/img/horse.png"/>
      <a href="https://site.com/about">About</a>
      <a href="https://site.com/about">About</a>
    </body>
  </html>`
)

func TestReturnsAListOfPageLinksFromHttp(t *testing.T) {
  site := "https://site.com"
  response := twoUrls
  extractor := LinkExtractor{&TestClient{response}, site}
  links := extractor.getLinks().PageLinks
  expectedLinks := []string{"https://site.com/about", "https://site.com/blog"}
  if !reflect.DeepEqual(links, expectedLinks) {
    t.Error("Expected: ", expectedLinks, ". Got: ", links)
  }
}

func TestReturnsAListOfAssetLinksFromHttp(t *testing.T) {
  site := "https://site.com"
  response := assetUrls
  extractor := LinkExtractor{&TestClient{response}, site}
  links := extractor.getLinks().ResourceLinks
  expectedLinks := []string{"https://site.com/myscript.js","https://site.com/mystyle.css", "https://site.com/img/horse.png"}
  if !reflect.DeepEqual(links, expectedLinks) {
    t.Error("Expected: ", expectedLinks, ". Got: ", links)
  }
}

func TestFiltersOutDuplicateUrls(t *testing.T) {
  site := "https://site.com"
  response := duplicates
  extractor := LinkExtractor{&TestClient{response}, site}
  pageLinks := extractor.getLinks().PageLinks
  assetLinks := extractor.getLinks().ResourceLinks
  expectedPageLinks := []string{"https://site.com/about"}
  expectedAssetLinks := []string{"https://site.com/myscript.js", "https://site.com/img/horse.png"}
  if !reflect.DeepEqual(pageLinks, expectedPageLinks) {
    t.Error("Expected: ", expectedPageLinks, ". Got: ", pageLinks)
  }

  if !reflect.DeepEqual(assetLinks, expectedAssetLinks) {
    t.Error("Expected: ", expectedAssetLinks, ". Got: ", assetLinks)
  }
}
