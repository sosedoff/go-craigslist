package craigslist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	reSiteId     = regexp.MustCompile(`https://(.*).craigslist.org`)
	reListingId  = regexp.MustCompile(`([\d]+).html`)
	reAttributes = regexp.MustCompile(`^(.+):\s?(.*)`)
	reImage      = regexp.MustCompile(`images.craigslist.org\/(.*)_[\d]+x[\d]+.jpg`)

	timeDetailsLayout = "2006-01-02T15:04:05-0700"
	timeIndexLayoyt   = "2006-01-02 15:04"

	errInvalidListing = errors.New("Invalid listing")
)

type Listing struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	URL         string            `json:"url"`
	Price       int               `json:"price"`
	Images      []Image           `json:"images"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	Location    *LatLng           `json:"location,omitempty"`
	PostedAt    *time.Time        `json:"posted_at,omitempty"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Image struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

func GetListing(url string) (*Listing, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseListing(resp.Body)
}

func ParseListing(reader io.Reader) (*Listing, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	url := doc.Find("link").AttrOr("href", "")
	if url == "" {
		return nil, errInvalidListing
	}

	listing := &Listing{
		Id:          reListingId.FindStringSubmatch(url)[1],
		URL:         url,
		Title:       doc.Find("#titletextonly").Text(),
		Description: parseDescription(doc),
		Price:       parsePrice(doc),
		Attributes:  parseAttributes(doc),
		Location:    parseLocation(doc),
		Images:      parseImages(doc),
	}

	if posted, updated := parseTimestamps(doc); posted != nil {
		listing.PostedAt = posted
		listing.UpdatedAt = updated
	}

	return listing, nil
}

func (l *Listing) JSON() (string, error) {
	data, err := json.Marshal(l)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseDescription(doc *goquery.Document) string {
	block := doc.Find("#postingbody").First()
	if block.Length() == 0 {
		return ""
	}

	block.Find("div").Remove()
	block.Find("a").Remove()

	return strings.TrimSpace(block.Text())
}

func parsePrice(doc *goquery.Document) int {
	var price int
	fmt.Sscanf(doc.Find(".postingtitletext > .price").Text(), "$%d", &price)
	return price
}

func parseAttributes(doc *goquery.Document) map[string]string {
	attrs := map[string]string{}

	doc.Find("p.attrgroup > span").Each(func(i int, s *goquery.Selection) {
		matches := reAttributes.FindStringSubmatch(s.Text())
		if len(matches) > 0 {
			key := strings.ToLower(strings.Replace(strings.TrimSpace(matches[1]), " ", "_", -1))
			value := strings.TrimSpace(matches[2])
			attrs[key] = value
		}
	})

	return attrs
}

func parseLocation(doc *goquery.Document) *LatLng {
	// Find the geo coordinates if they are available
	if block := doc.Find("#map").First(); block.Length() > 0 {
		var lat, lng float64

		fmt.Sscanf(block.AttrOr("data-latitude", ""), "%f", &lat)
		fmt.Sscanf(block.AttrOr("data-longitude", ""), "%f", &lng)

		if lat != 0 && lng != 0 {
			return &LatLng{lat, lng}
		}
	}

	return nil
}

func parseTimestamp(layout, text string) (time.Time, error) {
	ts, err := time.Parse(layout, text)
	if err != nil {
		return time.Now(), err
	}
	return ts.UTC(), nil
}

func parseTimestamps(doc *goquery.Document) (*time.Time, *time.Time) {
	timestamps := doc.Find(".postinginfos time").Map(func(i int, s *goquery.Selection) string {
		return s.AttrOr("datetime", "")
	})
	if len(timestamps) == 0 {
		return nil, nil
	}

	postedAt, err := parseTimestamp(timeDetailsLayout, timestamps[0])
	if err != nil {
		log.Println("cant parse time:", timestamps[0], err)
		return nil, nil
	}

	updatedAt := postedAt
	if len(timestamps) > 1 {
		updatedAt, _ = parseTimestamp(timeDetailsLayout, timestamps[1])
	}

	return &postedAt, &updatedAt
}

func parseImages(doc *goquery.Document) []Image {
	images := []Image{}
	thumbs := doc.Find("a.thumb")

	if thumbs.Length() > 0 {
		thumbs.Each(func(i int, s *goquery.Selection) {
			if href, exists := s.Attr("href"); exists {
				imageId := reImage.FindStringSubmatch(href)[1]

				images = append(images, Image{
					Small:  fmt.Sprintf("https://images.craigslist.org/%s_300x300.jpg", imageId),
					Medium: fmt.Sprintf("https://images.craigslist.org/%s_600x450.jpg", imageId),
					Large:  fmt.Sprintf("https://images.craigslist.org/%s_1200x900.jpg", imageId),
				})
			}
		})
	} else {
		src := doc.Find(".slide img").First().AttrOr("src", "")
		if src != "" {
			imageId := reImage.FindStringSubmatch(src)[1]

			images = append(images, Image{
				Small:  fmt.Sprintf("https://images.craigslist.org/%s_300x300.jpg", imageId),
				Medium: fmt.Sprintf("https://images.craigslist.org/%s_600x450.jpg", imageId),
				Large:  fmt.Sprintf("https://images.craigslist.org/%s_1200x900.jpg", imageId),
			})
		}
	}

	return images
}
