package feeder

import "testing"

func TestFeed(t *testing.T) {
	urilist := []string{
		//"http://localhost:8081/craigslist.rss",
		//"http://store.steampowered.com/feeds/news.xml", // This feed violates the rss spec.
		"http://cyber.law.harvard.edu/rss/examples/sampleRss091.xml",
		"http://cyber.law.harvard.edu/rss/examples/sampleRss092.xml",
		"http://cyber.law.harvard.edu/rss/examples/rss2sample.xml",
		"http://blog.case.edu/news/feed.atom",
	}

	var feed *Feed
	var err error

	for _, uri := range urilist {
		feed = New(5, true, chanHandler, itemHandler)

		if err = feed.Fetch(uri); err != nil {
			t.Errorf("%s >>> %s", uri, err)
			return
		}
	}
}

func chanHandler(feed *Feed, newchannels []*Channel) {
	println(len(newchannels), "new channel(s) in", feed.Url)
}

func itemHandler(feed *Feed, ch *Channel, newitems []*Item) {
	println(len(newitems), "new item(s) in", ch.Title, "of", feed.Url)
}
