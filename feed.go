/*
 Author: jim teeuwen <jimteeuwen@gmail.com>
 Dependencies: go-pkg-xmlx (http://github.com/jteeuwen/go-pkg-xmlx)

 This package allows us to fetch Rss and Atom feeds from the internet.
 They are parsed into an object tree which is a hyvrid of both the RSS and Atom
 standards.

 Supported feeds are:
 	- Rss v0.91, 0.91 and 2.0
 	- Atom 1.0

 The package allows us to maintain cache timeout management. This prevents us
 from querying the servers for feed updates too often and risk ip bams. Appart
 from setting a cache timeout manually, the package also optionally adheres to
 the TTL, SkipDays and SkipHours values specied in the feeds themselves.

 Note that the TTL, SkipDays and SkipHour fields are only part of the RSS spec.
 For Atom feeds, we use the CacheTimeout in the Feed struct.

 Because the object structure is a hybrid between both RSS and Atom specs, not
 all fields will be filled when requesting either an RSS or Atom feed. I have
 tried to create as many shared fields as possiblem but some of them simply do
 not occur in either the RSS or Atom spec.
*/
package feeder

import "os"
import "time"
import xmlx "github.com/jteeuwen/go-pkg-xmlx"
import "fmt"
import "strconv"
import "strings"

type ChannelHandler func(f *Feed, newchannels []*Channel)
type ItemHandler func(f *Feed, ch *Channel, newitems []*Item)

type Feed struct {
	// Custom cache timeout in minutes.
	CacheTimeout int

	// Make sure we adhere to the cache timeout specified in the feed. If
	// our CacheTimeout is higher than that, we will use that instead.
	EnforceCacheLimit bool

	// Type of feed. Rss, Atom, etc
	Type string

	// Version of the feed. Major and Minor.
	Version [2]int

	// Channels with content.
	Channels []*Channel

	// Url from which this feed was created.
	Url string

	// A notification function, used to notify the host when a new channel
	// has been found.
	chanhandler ChannelHandler

	// A notification function, used to notify the host when a new item
	// has been found for a given channel.
	itemhandler ItemHandler

	// Last time content was fetched. Used in conjunction with CacheTimeout
	// to ensure we don't get content too often.
	lastupdate int64
}

func New(cachetimeout int, enforcecachelimit bool, ch ChannelHandler, ih ItemHandler) *Feed {
	v := new(Feed)
	v.CacheTimeout = cachetimeout
	v.EnforceCacheLimit = enforcecachelimit
	v.Type = "none"
	v.chanhandler = ch
	v.itemhandler = ih
	return v
}

// This returns a timestamp of the last time the feed was updated.
// The value is in seconds.
func (this *Feed) LastUpdate() int64 { return this.lastupdate }

func (this *Feed) Fetch(uri string) (err os.Error) {
	if !this.CanUpdate() {
		return
	}

	this.Url = uri

	// Extract type and version of the feed so we can have the appropriate
	// function parse it (rss 0.91, rss 0.92, rss 2, atom etc).
	doc := xmlx.New()
	if err = doc.LoadUri(uri); err != nil {
		return
	}
	this.Type, this.Version = this.GetVersionInfo(doc)

	if ok := this.testVersions(); !ok {
		err = os.NewError(fmt.Sprintf("Unsupported feed: %s, version: %+v", this.Type, this.Version))
		return
	}

	chancount := len(this.Channels)
	if err = this.buildFeed(doc); err != nil || len(this.Channels) == 0 {
		return
	}

	// Notify host of new channels
	if chancount != len(this.Channels) && this.chanhandler != nil {
		this.chanhandler(this, this.Channels[chancount:])
	}

	// reset cache timeout values according to feed specified values (TTL)
	if this.EnforceCacheLimit && this.CacheTimeout < this.Channels[0].TTL {
		this.CacheTimeout = this.Channels[0].TTL
	}

	return
}

// This function returns true or false, depending on whether the CacheTimeout
// value has expired or not. Additionally, it will ensure that we adhere to the
// RSS spec's SkipDays and SkipHours values (if Feed.EnforceCacheLimit is set to
// true). If this function returns true, you can be sure that a fresh feed
// update will be performed.
func (this *Feed) CanUpdate() bool {
	// Make sure we are not within the specified cache-limit.
	// This ensures we don't request data too often.
	utc := time.UTC()
	if utc.Seconds()-this.lastupdate < int64(this.CacheTimeout*60) {
		return false
	}

	// If skipDays or skipHours are set in the RSS feed, use these to see if
	// we can update.
	if len(this.Channels) == 0 && this.Type == "rss" {
		if this.EnforceCacheLimit && len(this.Channels[0].SkipDays) > 0 {
			for _, v := range this.Channels[0].SkipDays {
				if v == utc.Weekday {
					return false
				}
			}
		}

		if this.EnforceCacheLimit && len(this.Channels[0].SkipHours) > 0 {
			for _, v := range this.Channels[0].SkipHours {
				if v == utc.Hour {
					return false
				}
			}
		}
	}

	this.lastupdate = utc.Seconds()
	return true
}

func (this *Feed) buildFeed(doc *xmlx.Document) (err os.Error) {
	switch this.Type {
	case "rss":
		err = this.readRss2(doc)
	case "atom":
		err = this.readAtom(doc)
	}
	return
}

func (this *Feed) testVersions() bool {
	switch this.Type {
	case "rss":
		if this.Version[0] > 2 || (this.Version[0] == 2 && this.Version[1] > 0) {
			return false
		}

	case "atom":
		if this.Version[0] > 1 || (this.Version[0] == 1 && this.Version[1] > 0) {
			return false
		}

	default:
		return false
	}

	return true
}

func (this *Feed) GetVersionInfo(doc *xmlx.Document) (ftype string, fversion [2]int) {
	node := doc.SelectNode("http://www.w3.org/2005/Atom", "feed")
	if node == nil {
		goto rss
	}
	ftype = "atom"
	fversion = [2]int{1, 0}
	return

rss:
	node = doc.SelectNode("", "rss")
	if node == nil {
		goto end
	}
	ftype = "rss"
	version := node.GetAttr("", "version")
	p := strings.Index(version, ".")
	major, _ := strconv.Atoi(version[0:p])
	minor, _ := strconv.Atoi(version[p+1 : len(version)])
	fversion = [2]int{major, minor}
	return

end:
	ftype = "unknown"
	fversion = [2]int{0, 0}
	return
}
