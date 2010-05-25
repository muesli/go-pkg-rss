package feeder

import "testing"

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
		feed = New(5, true)

		if err = feed.Fetch(uri); err != nil {
			t.Errorf("%s >>> %s", uri, err)
		}
	}
}
