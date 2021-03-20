package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mattn/go-mastodon"
)

var (
	mastodonServer       string
	mastodonClientID     string
	mastodonClientSecret string
	mastodonAccessToken  string
)

func init() {
	flag.StringVar(&mastodonServer, "mastodon-server", "", "Mastodon server")
	flag.StringVar(&mastodonClientID, "mastodon-key", "", "Mastodon key")
	flag.StringVar(&mastodonClientSecret, "mastodon-secret", "", "Mastodon secret")
	flag.StringVar(&mastodonAccessToken, "mastodon-access-token", "", "Mastodon access token")
}

func main() {
	flag.Parse()

	c := mastodon.NewClient(&mastodon.Config{
		Server:       mastodonServer,
		ClientID:     mastodonClientID,
		ClientSecret: mastodonClientSecret,
		AccessToken:  mastodonAccessToken,
	})
	a, err := c.GetAccountCurrentUser(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	var pg mastodon.Pagination
	var statuses []*mastodon.Status
	for {
		s, err := c.GetAccountStatuses(context.Background(), a.ID, &pg)
		if err != nil {
			log.Fatalln(err)
		}

		statuses = append(statuses, s...)
		if pg.MaxID == "" || pg.MinID == "" {
			break
		}
		pg.MinID = ""
		pg.SinceID = ""

		time.Sleep(2 * time.Second)
	}

	for _, s := range statuses {
		fmt.Printf("%v %s\n\n", s.CreatedAt, s.Content)

		if err := c.DeleteStatus(context.Background(), s.ID); err != nil {
			log.Fatalln(err)
		}

		time.Sleep(2 * time.Second)
	}
}
