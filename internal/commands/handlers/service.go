package handlers

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/config"
	"github.com/git-cst/bootdev_gator/internal/database"
)

func HandlerAgg(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("expected to receive time duration string as command: %s", cmd.Args)
	}
	durationString := cmd.Args[0]

	timeBetweenReqs, err := time.ParseDuration(durationString)
	if err != nil {
		return fmt.Errorf("error parsing %s\nerror was %v", timeBetweenReqs, err)
	}

	fmt.Println("==================================================")
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
		scrapeFeeds(s)
	}
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

func scrapeFeeds(s *config.State) error {
	ctx := context.Background()
	feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		// Handle no rows error
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no feeds registered in database")
		}
		// Handle other errors
		return err
	}

	params := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: feed.ID,
	}
	fetchedFeed, err := s.Db.MarkFeedFetched(ctx, params)
	if err != nil {
		// Handle no rows error
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no feed returned by mark feed fetched")
		}
		// Handle other errors
		return err
	}

	rssFeed, err := fetchFeed(s.Config, ctx, fetchedFeed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Fetched the following items from rssfeed %v:\n\n", fetchedFeed.Name)
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
	fmt.Println("==================================================")

	return nil
}
