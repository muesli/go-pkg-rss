package feeder

import (
	"io/ioutil"
	"testing"
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

func Test_NewItem(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/initial.atom")
	feed := New(1, true, chanHandler, itemHandler)
	err := feed.FetchBytes("http://example.com", content, nil)
	if err != nil {
		t.Error(err)
	}

	content, _ = ioutil.ReadFile("testdata/initial_plus_one_new.atom")
	feed.FetchBytes("http://example.com", content, nil)
	expected := "Second title"
	if len(items) != 1 {
		t.Errorf("Expected %s new item, got %s", 1, len(items))
	}

	if expected != items[0].Title {
		t.Errorf("Expected %s, got %s", expected, items[0].Title)
	}
}

func Test_AtomAuthor(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/idownload.atom")
	if err != nil {
		t.Errorf("unable to load file")
	}
	feed := New(1, true, chanHandler, itemHandler)
	err = feed.FetchBytes("http://example.com", content, nil)

	item := feed.Channels[0].Items[0]
	expected := "Cody Lee"
	if item.Author.Name != expected {
		t.Errorf("Expected author to be %s but found %s", expected, item.Author.Name)
	}
}

func Test_RssAuthor(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/boing.rss")
	feed := New(1, true, chanHandler, itemHandler)
	feed.FetchBytes("http://example.com", content, nil)

	item := feed.Channels[0].Items[0]
	expected := "Cory Doctorow"
	if item.Author.Name != expected {
		t.Errorf("Expected author to be %s but found %s", expected, item.Author.Name)
	}
}

func Test_CData(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/iosBoardGameGeek.rss")
	feed := New(1, true, chanHandler, itemHandler)
	feed.FetchBytes("http://example.com", content, nil)

	item := feed.Channels[0].Items[0]
	expected := `<p>abc<div>"def"</div>ghi`
	if item.Description != expected {
		t.Errorf("Expected item.Description to be [%s] but item.Description=[%s]", expected, item.Description)
	}
}

func Test_Link(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/ignoredLink.rss")
	feed := New(1, true, chanHandler, itemHandler)
	feed.FetchBytes("http://example.com", content, nil)

	channel := feed.Channels[0]
	item := channel.Items[0]

	channelLinkExpected := "http://www.conservatives.com/XMLGateway/RSS/News.xml"
	itemLinkExpected := "http://www.conservatives.com/News/News_stories/2013/09/Dr_Tania_Mathias_chosen_to_stand_up_for_local_people_in_Twickenham.aspx"

	if channel.Links[0].Href != channelLinkExpected {
		t.Errorf("Expected author to be %s but found %s", channelLinkExpected, channel.Links[0].Href)
	}

	if item.Links[0].Href != itemLinkExpected {
		t.Errorf("Expected author to be %s but found %s", itemLinkExpected, item.Links[0].Href)
	}
}

func chanHandler(feed *Feed, newchannels []*Channel) {
	println(len(newchannels), "new channel(s) in", feed.Url)
}

func itemHandler(feed *Feed, ch *Channel, newitems []*Item) {
	items = newitems
	println(len(newitems), "new item(s) in", ch.Title, "of", feed.Url)
}
