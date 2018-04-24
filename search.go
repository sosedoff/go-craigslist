package craigslist

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SearchResults struct {
	RangeFrom  int       `json:"range_from"`
	RangeTo    int       `json:"range_to"`
	TotalCount int       `json:"total_count"`
	Listings   []Listing `json:"listings"`
}

func Search(siteId string, opts SearchOptions) (*SearchResults, error) {
	url := fmt.Sprintf("https://%s.craigslist.org/search/%s", siteId, opts.Category)

	req, err := http.NewRequest("get", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range opts.Params {
		q.Add(k, v)
	}
	if opts.Skip > 0 {
		q.Add("s", strconv.Itoa(opts.Skip))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSearchResults(resp.Body)
}

func ParseSearchResults(reader io.Reader) (*SearchResults, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	results := &SearchResults{}

	fmt.Sscanf(doc.Find(".rangeFrom").First().Text(), "%d", &results.RangeFrom)
	fmt.Sscanf(doc.Find(".rangeTo").First().Text(), "%d", &results.RangeTo)
	fmt.Sscanf(doc.Find(".totalcount").First().Text(), "%d", &results.TotalCount)

	listings := []Listing{}

	doc.Find("li.result-row").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".result-title").First()
		gallery := s.Find(".result-image").First().AttrOr("data-ids", "")
		priceText := s.Find(".result-price").First().Text()
		price := 0
		images := []Image{}

		fmt.Sscanf(priceText, "$%d", &price)
		if price < 0 {
			price = 0
		}

		if gallery != "" {
			for _, str := range strings.Split(gallery, ",") {
				imageId := strings.Split(str, ":")[1]

				images = append(images, Image{
					Small:  fmt.Sprintf("https://images.craigslist.org/%s_300x300.jpg", imageId),
					Medium: fmt.Sprintf("https://images.craigslist.org/%s_600x450.jpg", imageId),
					Large:  fmt.Sprintf("https://images.craigslist.org/%s_1200x900.jpg", imageId),
				})
			}
		}

		listing := Listing{
			Id:     s.AttrOr("data-pid", ""),
			Title:  title.Text(),
			URL:    title.AttrOr("href", ""),
			Price:  price,
			Images: images,
		}

		if dt, exists := s.Find("time.result-date").First().Attr("datetime"); exists {
			if postedAt, err := parseTimestamp(timeIndexLayoyt, dt); err == nil {
				listing.PostedAt = &postedAt
				listing.UpdatedAt = &postedAt
			}
		}

		listings = append(listings, listing)
	})

	results.Listings = listings

	return results, nil
}
