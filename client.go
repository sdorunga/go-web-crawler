package main

import (
  "net/http"
  "io/ioutil"
  "fmt"
)
type HttpClient interface {
  Get(string) string
}

type Client struct {
  client http.Client
}

func (this *Client) Get(url string) string{
  resp, err := this.client.Get(url)
  if err != nil {
    fmt.Println("Error fetching url:", url)
    return ""
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Could not read body for request:", url)
    return ""
  }
  return string(body)
}

func NewClient() HttpClient {
  return &Client{http.Client{}}
}
