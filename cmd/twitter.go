package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var accessKey string
var accessSecret string
var twitterKey string
var twitterSecret string

func init() {
	flag.StringVar(&accessKey, "access-key", "", "Twitter access key")
	flag.StringVar(&accessSecret, "access-secret", "", "Twitter access secret")
	flag.StringVar(&twitterKey, "twitter-key", "", "Twitter consumer key")
	flag.StringVar(&twitterSecret, "twitter-secret", "", "Twitter consumer secret")
}

func main() {
	flag.Parse()

	config := oauth1.NewConfig(twitterKey, twitterSecret)
	token := oauth1.NewToken(accessKey, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		Count: 1000,
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Deleting tweets ...")
	for _, t := range tweets {
		_, _, err = client.Statuses.Destroy(t.ID, &twitter.StatusDestroyParams{})
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("Deleting likes ...")
	likes, _, err := client.Favorites.List(&twitter.FavoriteListParams{
		Count: 10000,
	})
	if err != nil {
		log.Fatalln(err)
	}

	for _, l := range likes {
		_, _, err = client.Favorites.Destroy(&twitter.FavoriteDestroyParams{
			ID: l.ID,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("Stop all following ...")
	var cursor int64
	for {
		friends, _, err := client.Friends.List(&twitter.FriendListParams{
			Count:  10000,
			Cursor: cursor,
		})
		if err != nil {
			log.Fatalln(err)
		}

		for _, f := range friends.Users {
			_, _, err = client.Friendships.Destroy(&twitter.FriendshipDestroyParams{
				UserID: f.ID,
			})
			if err != nil {
				log.Fatalln(err)
			}
		}

		cursor = friends.NextCursor
		if cursor < 0 {
			break
		}
	}
}
