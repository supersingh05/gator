package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/supersingh05/gator/internal/database"
	"github.com/supersingh05/gator/internal/rss"
)

func handler_unfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("not enough args")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedForUser(context.Background(), database.DeleteFeedForUserParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Printf("URL: %s deletes for user: %s", feed.Url, user.Name.String)
	return nil
}

func handler_following(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedsForUser(context.Background(), s.config.CurrentUserId)
	if err != nil {
		return err
	}
	for _, f := range feeds {
		fmt.Printf("%s\n", f.Name)
	}
	return nil
}

func handler_followfeeds(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("Not enough args")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	_, err = s.db.CreateUserFeed(context.Background(), database.CreateUserFeedParams{UserID: s.config.CurrentUserId, FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Println("User is following feed now")

	return nil
}

func handler_feeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeedsWithUserName(context.Background())
	if err != nil {
		return err
	}
	for _, f := range feeds {
		fmt.Printf("Name: %s Url: %s CreatedBy: %s \n", f.Feedname, f.Url, f.Username.String)
	}
	return nil
}

func handler_addfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("Not enough args")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	ctx := context.Background()
	cfd := database.CreateFeedParams{Name: name, Url: url, CreatedBy: s.config.CurrentUserId}
	tx, err := s.dbObj.BeginTx(ctx, nil)
	if err != nil {
		return errors.New("couldnt make tx")
	}

	defer tx.Rollback()
	qtx := s.db.WithTx(tx)
	dbf, err := qtx.CreateFeed(ctx, cfd)
	if err != nil {
		return err
	}
	cfp := database.CreateUserFeedParams{UserID: s.config.CurrentUserId, FeedID: dbf.ID}
	_, err = qtx.CreateUserFeed(ctx, cfp)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("couldnt finish txn")
		return err
	}
	fmt.Printf("%v/n", dbf)
	return nil
}

func handler_browse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		limit, _ = strconv.Atoi(cmd.args[0])
	}
	ctx := context.Background()
	posts, err := s.db.GetPostsForUser(ctx, s.config.CurrentUserId)
	if err != nil {
		return err
	}
	for _, p := range posts[:limit] {
		fmt.Printf("%v \n", p)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("Not enough args")
	}
	timeBetweenRequests := cmd.args[0]
	t, err := time.ParseDuration(timeBetweenRequests)
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s \n", timeBetweenRequests)
	ticker := time.NewTicker(t)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	ctx := context.Background()
	tx, err := s.dbObj.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()
	qtx := s.db.WithTx(tx)
	feed, err := qtx.GetNextFeedToFetch(ctx)
	if err != nil {
		return
	}
	t := time.Now()
	var lfa sql.NullTime
	lfa.Scan(t)
	_, err = qtx.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{UpdatedAt: t, LastFetchedAt: lfa, ID: feed.ID})
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("couldnt finish txn")
		return
	}
	feeds, err := fetchFeed(ctx, feed.Url)
	for _, f := range feeds.Channel.Item {
		t, err := parseAnyFormat(f.PubDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nt := sql.NullTime{Time: t, Valid: true}
		p := database.CreatePostParams{PublishedAt: nt, Title: f.Title, Description: f.Description, Url: f.Link, CreatedBy: s.config.CurrentUserId, FeedID: feed.ID}
		post, err := s.db.CreatePost(ctx, p)
		if err != nil {
			if isUniqueViolation(err) {
				fmt.Println("had unique violation, but continuing")
				continue
			}
			fmt.Printf("error with saving, returning %v \n", err)
			return
		}
		fmt.Printf("Post: %v %v for feed: %v saved \n", post.ID, post.Title, feed.Name)
	}

}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == "23505" // unique_violation
	}
	return false
}

func parseAnyFormat(dateString string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC3339,
		time.ANSIC,
		"02/01/2006 15:04:05",
		"2 Jan 06 03:04PM",
		// Add more possible formats here
	}

	for _, format := range formats {
		t, err := time.Parse(format, dateString)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse time with any known format: %s", dateString)
}

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var feed rss.RSSFeed
	defer res.Body.Close()
	decoder := xml.NewDecoder(res.Body)
	if err := decoder.Decode(&feed); err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for idx := range feed.Channel.Item {
		feed.Channel.Item[idx].Description = html.UnescapeString(feed.Channel.Item[idx].Description)
		feed.Channel.Item[idx].Title = html.UnescapeString(feed.Channel.Item[idx].Title)
	}

	return &feed, nil
}
