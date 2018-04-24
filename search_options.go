package craigslist

import (
	"net/url"
	"strconv"
)

type SearchOptions struct {
	Category         string
	Query            string
	TitlesOnly       bool
	HasImage         bool
	PostedToday      bool
	BundleDuplicates bool
	IncludeNearby    bool
	MinPrice         int
	MaxPrice         int
	Params           map[string]string
	Skip             int
}

func (opts SearchOptions) query() url.Values {
	q := url.Values{}

	if opts.Query != "" {
		q.Add("query", opts.Query)
	}
	if opts.TitlesOnly {
		q.Add("srchType", "T")
	}
	if opts.HasImage {
		q.Add("hasPic", "1")
	}
	if opts.PostedToday {
		q.Add("postedToday", "1")
	}
	if opts.BundleDuplicates {
		q.Add("bundleDuplicates", "1")
	}
	if opts.IncludeNearby {
		q.Add("searchNearby", "1")
	}
	if opts.MinPrice > 0 {
		q.Add("min_price", strconv.Itoa(opts.MinPrice))
	}
	if opts.MaxPrice > 0 {
		q.Add("max_price", strconv.Itoa(opts.MaxPrice))
	}

	// Add any custom params that are not common on the search page
	for k, v := range opts.Params {
		q.Add(k, v)
	}

	return q
}
