package feeder

import "testing"
import "os"

func TestFeed(t *testing.T) {
	urilist := []string{
		"http://cyber.law.harvard.edu/rss/examples/sampleRss091.xml",
		"http://cyber.law.harvard.edu/rss/examples/sampleRss092.xml",
		"http://cyber.law.harvard.edu/rss/examples/rss2sample.xml",
		"http://blog.case.edu/news/feed.atom",
	}

	var feed *Feed
	var err os.Error

	for _, uri := range urilist {
		feed = New(5, true, chanHandler, itemHandler)

		if err = feed.Fetch(uri); err != nil {
			t.Errorf("%s >>> %s", uri, err)
		}
	}

	/*
		Output of handlers:

		6 new item(s) in WriteTheWeb of http://cyber.law.harvard.edu/rss/examples/sampleRss091.xml
		1 new channel(s) in http://cyber.law.harvard.edu/rss/examples/sampleRss091.xml
		21 new item(s) in Dave Winer: Grateful Dead of http://cyber.law.harvard.edu/rss/examples/sampleRss092.xml
		1 new channel(s) in http://cyber.law.harvard.edu/rss/examples/sampleRss092.xml
		4 new item(s) in Liftoff News of http://cyber.law.harvard.edu/rss/examples/rss2sample.xml
		1 new channel(s) in http://cyber.law.harvard.edu/rss/examples/rss2sample.xml
		15 new item(s) in Blog@Case of http://blog.case.edu/news/feed.atom
		1 new channel(s) in http://blog.case.edu/news/feed.atom
	*/
}

func chanHandler(feed *Feed, newchannels []*Channel) {
	//println(len(newchannels), "new channel(s) in", feed.Url)
}

func itemHandler(feed *Feed, ch *Channel, newitems []*Item) {
	//println(len(newitems), "new item(s) in", ch.Title, "of", feed.Url)
}
