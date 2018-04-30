package craigslist

import (
	"net/url"
	"strconv"
)

type SearchOptions struct {
	Category         string `form:"category"`
	Query            string `form:"query"`
	TitlesOnly       bool   `form:"titles_only"`
	HasImage         bool   `form:"has_image"`
	PostedToday      bool   `form:"posted_today"`
	BundleDuplicates bool   `form:"bundle_duplicates"`
	IncludeNearby    bool   `form:"include_nearby"`
	MinPrice         int    `form:"min_price"`
	MaxPrice         int    `form:"max_price"`
	Skip             int    `form:"skip"`
	Params           map[string]string
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
	if opts.Skip > 0 {
		q.Add("s", strconv.Itoa(opts.Skip))
	}

	// Add any custom params that are not common on the search page
	for k, v := range opts.Params {
		q.Add(k, v)
	}

	return q
}
