package main

/*
This is a minimal sample application, demonstrating how to set up an RSS feed
for regular polling of new channels/items.
*/

import (
	"fmt"
	"os"
	"time"
	rss "github.com/jteeuwen/go-pkg-rss"
)

func main() {
	// This sets up a new feed and polls it for new channels/items in
	// a separate goroutine. Invoke it with 'go PollFeed(..)' to have the
	// polling performed in a separate goroutine, so you can continue with
	// the rest of your program.
	PollFeed("http://cyber.law.harvard.edu/rss/examples/sampleRss091.xml", 5)
}

func PollFeed(uri string, timeout int) {
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		if err := feed.Fetch(uri); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		<-time.After(feed.SecondsTillUpdate() * 1e9)
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	fmt.Printf("%d new item(s) in %s\n", len(newitems), feed.Url)
}
