package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig("consumerKey", "consumerSecret")
	token := oauth1.NewToken("accessToken", "accessSecret")

	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	tweet, resp, err := client.Statuses.Update("▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░ 70.8%", nil)
	fmt.Println(err)
	fmt.Println(resp)
	fmt.Println(tweet)
}
