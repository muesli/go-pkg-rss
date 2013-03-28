package feeder

import (
    "testing"
	"io/ioutil"
)

var items []*Item

func TestFeed(t *testing.T) {
	urilist := []string{
		//"http://cyber.law.harvard.edu/rss/examples/sampleRss091.xml", // Non-utf8 encoding.
		"http://store.steampowered.com/feeds/news.xml", // This feed violates the rss spec.
		"http://cyber.law.harvard.edu/rss/examples/sampleRss092.xml",
		"http://cyber.law.harvard.edu/rss/examples/rss2sample.xml",
		"http://blog.case.edu/news/feed.atom",
	}

	var feed *Feed
	var err error

	for _, uri := range urilist {
		feed = New(5, true, chanHandler, itemHandler)

		if err = feed.Fetch(uri, nil); err != nil {
			t.Errorf("%s >>> %s", uri, err)
			return
		}
	}
}

func Test_AtomAuthor(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/idownload.atom")
    if err != nil {
		t.Errorf("unable to load file")
    }
	feed := New(1, true, chanHandler, itemHandler)
	err = feed.FetchBytes("http://example.com", content, nil)

	item := items[0]
	expected := "Cody Lee"
	if item.Author.Name != expected {
		t.Errorf("Expected author to be %s but found %s", expected, item.Author.Name)
	}
}

func Test_RssAuthor(t *testing.T) {
    content, _ := ioutil.ReadFile("testdata/boing.rss")
    feed := New(1, true, chanHandler, itemHandler)
    feed.FetchBytes("http://example.com", content, nil)

    item := items[0]
    expected := "Cory Doctorow"
    if item.Author.Name != expected {
        t.Errorf("Expected author to be %s but found %s", expected, item.Author.Name)
    }
}

func chanHandler(feed *Feed, newchannels []*Channel) {
	println(len(newchannels), "new channel(s) in", feed.Url)
}

func itemHandler(feed *Feed, ch *Channel, newitems []*Item) {
    items = newitems
	println(len(newitems), "new item(s) in", ch.Title, "of", feed.Url)
}
