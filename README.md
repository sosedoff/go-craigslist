# go-craigslist

[WIP] Craigslist.org wrapper for Go

## Install

```
go get -u github.com/sosedoff/go-craigslist
```

## Usage

```golang
package main

import (
  "log"

  "github.com/sosedoff/go-craigslist"
)

func main() {
  opts := craigslist.SearchOptions{
    Category: "cto", // cars+trucks
    Params: map[string]string{
      "query": "honda",
    },
  }

  results, err := craigslist.Search("chicago", opts)
  if err != nil {
    log.Fatal(err)
  }

  log.Println(results)
}
```