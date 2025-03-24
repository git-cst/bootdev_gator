package handlers

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/config"
)

func HandlerAgg(s *config.State, cmd commands.Command) error {
	/*if len(cmd.Args) < 1 {
		return fmt.Errorf("no feed passed to the aggregate handler: %v", cmd.Args)
	}
	feedUrl := cmd.Args[1]*/
	feedUrl := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()

	rssFeed, err := fetchFeed(s.Config, ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("error whilst fetching feed: %v", err)
	}

	jsonData, err := json.MarshalIndent(rssFeed, "", "	")
	if err != nil {
		return fmt.Errorf("error whilst prettifying rssfeed data: %v", err)
	}

	fmt.Println(string(jsonData))

	return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(c *config.Config, ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feed RSSFeed
	xml.Unmarshal(data, &feed)

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
