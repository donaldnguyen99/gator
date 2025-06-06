package cli

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/donaldnguyen99/gator/internal/database"
	"github.com/donaldnguyen99/gator/internal/rss"
	"github.com/google/uuid"
)

type handler func(*state, command) error

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires 1 argument (username)")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error handling login: %w", err)
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error handling login: %w", err)
	}
	fmt.Printf("The user %s has been set\n", cmd.args[0])
	return nil
}

func handlerRegisterUser(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires 1 argument (username)")
	}
	u, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("error handling register user: %w", err)
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error handling register user: %w", err)
	}
	fmt.Printf("The user %v has been created and set\n", u)
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("get users requires 0 arguments")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error handling get users: %w", err)
	}

	for _, u := range users {
		if u.Name == s.config.CurrentUserName {
			fmt.Printf("* %v (current)\n", u.Name)
		} else {
			fmt.Printf("* %v\n", u.Name)
		}
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset requires 0 arguments")
	}
	if err := s.config.SetUser(""); err != nil {
		return fmt.Errorf("error handling reset: %w", err)
	}
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("error handling reset: %w", err)
	}
	fmt.Println("The database has been reset")
	return nil
}

func handlerAggregateFeeds(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg requires 1 argument, time_betweeen_reqs")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	timeBetweenRequests = min(max(timeBetweenRequests, 1*time.Second), 24*time.Hour)
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests.String())

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("addfeed requires 2 arguments (name, url), got %v arguments", len(cmd.args))
	}

	name, url := cmd.args[0], cmd.args[1]

	if s.config.CurrentUserName == "" {
		return fmt.Errorf("no user logged in")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error handling add feed: %w", err)
	}
	fmt.Printf("Feed %v, %v has been added\n", feed.Name, feed.Url)

	feed_follows_row, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:       uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:   user.ID,
		FeedID:   feed.ID,
	})
	if err != nil || feed_follows_row.UserID != user.ID || feed_follows_row.FeedID != feed.ID {
		return fmt.Errorf("error creating feed following while adding feed: %w", err)
	}

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsWithUserName(context.Background())
	if err != nil {
		return fmt.Errorf("error handling list feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Feed %v, %v, %v\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow requires 1 argument, feed_url")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("no feed found: %w", err)
	}

	feed_follows_row, err := s.db.CreateFeedFollow(
		context.Background(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil || feed_follows_row.FeedName != feed.Name || feed_follows_row.UserName != user.Name {
		return fmt.Errorf("error handling follow feed: %w", err)
	}

	fmt.Printf("Feed %v, has been followed by %v\n", feed.Name, user.Name)

	return nil
}

func handlerListFollows(s *state, cmd command, user database.User) error {
	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error handling list follows: %w", err)
	}
	for _, feed_follow := range feed_follows {
		fmt.Printf("Feed %v, %v\n", feed_follow.FeedName, feed_follow.UserName)
	}
	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unfollow requires 1 argument, feed_url")
	}
	err := s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserName: user.Name,
		FeedUrl:  cmd.args[0],
	})
	if err != nil {
		return err
	}
	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	// rssFeed.Print()
	for _, item := range rssFeed.Channel.Items {

		pubDate, err := item.ParsePubDate()
		if err != nil {
			return err
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      nextFeed.ID,
		})

		if err != nil && !strings.Contains(err.Error(), "duplicate key") {
			log.Print(err)
		}
	}

	return nil
}

func handlerBrowsePosts(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		res, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = res
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Name:  user.Name,
		Limit: int32(limit),
	})

	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Description: %s\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Printf("Published: %s\n\n", post.PublishedAt.Format(time.RFC1123))
	}

	return nil
}