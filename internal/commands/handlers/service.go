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
	"github.com/lib/pq"
)

func HandlerAgg(s *config.State, cmd commands.Command) error {
	s.LogInfo("Start aggregator service")
	if len(cmd.Args) < 1 {
		return fmt.Errorf("expected to receive time duration string as command: %s", cmd.Args)
	}
	durationString := cmd.Args[0]

	timeBetweenReqs, err := time.ParseDuration(durationString)
	if err != nil {
		s.LogError("Error parsing %s to time.Duration. Error was %v", timeBetweenReqs, err)
		return fmt.Errorf("error parsing %s\nerror was %v", timeBetweenReqs, err)
	}

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		s.LogInfo("Collecting feeds every %v", timeBetweenReqs)
		scrapeFeeds(s)
	}
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
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

	for i := range feed.Channel.Items {
		feed.Channel.Items[i].Title = html.UnescapeString(feed.Channel.Items[i].Title)
		feed.Channel.Items[i].Description = html.UnescapeString(feed.Channel.Items[i].Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *config.State) error {
	ctx := context.Background()
	feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		// Handle no rows error
		if errors.Is(err, sql.ErrNoRows) {
			s.LogInfo("No feeds registered in database")
			return fmt.Errorf("no feeds registered in database")
		}
		// Handle other errors
		s.LogError("Failed to get next feed: %v", err)
		return err
	}

	markFetchedParams := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: feed.ID,
	}
	fetchedFeed, err := s.Db.MarkFeedFetched(ctx, markFetchedParams)
	if err != nil {
		// Handle no rows error
		if errors.Is(err, sql.ErrNoRows) {
			s.LogError("No feed returned by mark feed fetched.")
			return fmt.Errorf("no feed returned by mark feed fetched")
		}
		// Handle other errors
		return err
	}

	rssFeed, err := fetchFeed(s.Config, ctx, fetchedFeed.Url)
	if err != nil {
		return err
	}

	s.LogInfo("Fetching feeds from %v", fetchedFeed.Name)
	for _, item := range rssFeed.Channel.Items {
		var descriptionNull sql.NullString
		if item.Description != "" {
			descriptionNull = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		} else {
			descriptionNull = sql.NullString{
				Valid: false,
			}
		}

		createPostParams := database.CreatePostParams{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: descriptionNull,
			PublishedAt: parseTimeString(item.PubDate),
			FeedID:      fetchedFeed.ID,
		}

		_, err := s.Db.CreatePost(ctx, createPostParams)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" { // PostgreSQL error code for unique violation
					// OK expected behaviour
					continue
				}
				// Log the error
				s.LogError("Could not create the post for fetched feed %s (%v). Item failed was %s", fetchedFeed.Name, fetchedFeed.ID, item.Title)
			}
		}
		s.LogInfo("Created post for fetched feed %v. Post title: %s (Published: %s)", fetchedFeed.Name, item.Title, item.PubDate)
	}

	s.LogInfo("Feed scraping completed successfully")
	return nil
}
