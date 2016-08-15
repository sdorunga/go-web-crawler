package main

import (
  "golang.org/x/net/html"
  "strings"
  "fmt"
  "net/url"
)

type LinkExtractor struct {
  client HttpClient
  url string
}

type SiteLinks struct {
  PageLinks []string
  ResourceLinks []string
}

func (this *LinkExtractor) getLinks() SiteLinks {
  body := this.client.Get(this.url)
  parsedBody, err := html.Parse(strings.NewReader(body))
  if err != nil {
    fmt.Println("Failed to parse the HTTP response")
  }
  pageLinks := this.collectPageLinks(parsedBody, []string{})
  resourceLinks := this.collectResourceLinks(parsedBody, []string{})
  return SiteLinks{pageLinks, resourceLinks}
}

func (this *LinkExtractor) collectPageLinks(rootNode *html.Node, links []string) []string {
  for node := rootNode.FirstChild; node != nil; node = node.NextSibling {
    if link := this.extractInternalLink(node); link != "" {
          links = append(links, link)
    }
    links = this.collectPageLinks(node, links)
  }

  return links
}

func (this *LinkExtractor) collectResourceLinks(rootNode *html.Node, links []string) []string {
  for node := rootNode.FirstChild; node != nil; node = node.NextSibling {
    if link := this.extractInternalResourceLink(node); link != "" {
          links = append(links, link)
    }
    links = this.collectResourceLinks(node, links)
  }

  return links
}

func (this *LinkExtractor) extractInternalResourceLink(node *html.Node) (link string) {
  if node.Type == html.ElementNode && (node.Data == "link" || node.Data == "script" || node.Data == "img" ) {
    for _, attribute := range node.Attr {
      if attribute.Key == "href" || attribute.Key == "src" {
        link = this.absolutifyUrl(attribute.Val).String()
      }
    }
  }
  return
}

func (this *LinkExtractor) extractInternalLink(node *html.Node) (link string) {
  if node.Type == html.ElementNode && node.Data == "a" {
    for _, attribute := range node.Attr {
      if attribute.Key == "href" {
        if absoluteLink := this.absolutifyUrl(attribute.Val); this.isInternalLink(absoluteLink) {
          link = absoluteLink.String()
        }
      }
    }
  }
  return
}

func (this *LinkExtractor) absolutifyUrl(link string) *url.URL{
  linkUrl, err := url.Parse(link)
  if err != nil {
    fmt.Println("Malformed Url", link)
  }
  if linkUrl.Host == "" {
    linkUrl.Scheme = "https"
    linkUrl.Host = "site.com"
  }
  return linkUrl
}

func (this *LinkExtractor) isInternalLink(link *url.URL) bool{
  parsedUrl, err := url.Parse(this.url)
  if err != nil {
    fmt.Println("Malformed Url", this.url)
  }
  host := parsedUrl.Host
  return link.Host == host
}
