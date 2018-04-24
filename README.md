# go-craigslist

[![Build Status](https://travis-ci.org/sosedoff/go-craigslist.svg)](https://travis-ci.org/sosedoff/go-craigslist)
[![GoDoc](https://godoc.org/github.com/sosedoff/go-craigslist?status.svg)](https://godoc.org/github.com/sosedoff/go-craigslist)

Craigslist.org wrapper for Go

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
    Query:    "honda",
    HasImage: true,
    MinPrice: 10000,
    MaxPrice: 20000,
  }

  // Perform a search
  result, err := craigslist.Search("chicago", opts)
  if err != nil {
    log.Fatal(err)
  }

  for _, listing := range result.Listings {
    log.Println(listing.JSON())
  }

  // Fetch listing details
  listing, err := craigslist.GetListing(result.Listings[0].URL)
  if err != nil {
    log.Fatal(err)
  }

  log.Println(listing)
}
```

## Test

```
make test
```

## License

MIT
