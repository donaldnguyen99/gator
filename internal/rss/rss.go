package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("User-Agent", "gator")

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch feed: status code %d", resp.StatusCode)
	}

	var rssFeed RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&rssFeed); err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range rssFeed.Channel.Items {
		rssFeed.Channel.Items[i].Title = html.UnescapeString(rssFeed.Channel.Items[i].Title)
		rssFeed.Channel.Items[i].Description = html.UnescapeString(rssFeed.Channel.Items[i].Description)
	}

	return &rssFeed, nil
}

func AggregateFeeds(ctx context.Context, feedUrls []string) ([]*RSSFeed, error) {
	var rssFeeds []*RSSFeed
	for _, feedUrl := range feedUrls {
		feed, err := FetchFeed(ctx, feedUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch feed %s: %w", feedUrl, err)
		}
		rssFeeds = append(rssFeeds, feed)
	}

	return rssFeeds, nil
}

func (r *RSSFeed) Print() {
	fmt.Printf("Channel Title: %s\n", r.Channel.Title)
	fmt.Printf("  Link: %s\n", r.Channel.Link)
	fmt.Printf("  Description: %s\n", r.Channel.Description)
	fmt.Print( "  Items:\n")
	r.PrintItem(true, false, false, false)
}

func (r *RSSFeed) PrintItem(title bool, link bool, desc bool, date bool) {
	for _, item := range r.Channel.Items {
		if title {
			fmt.Printf("    Title: %s\n", item.Title)
		}
		if link {
			fmt.Printf("    Link: %s\n", item.Link)
		}
		if desc {
			fmt.Printf("    Description: %s\n", item.Description)
		}
		if date {
			fmt.Printf("    PubDate: %s\n", item.PubDate)
		}
	}
}