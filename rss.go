package feeder

import (
	"errors"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"
)

func (this *Feed) readRss2(doc *xmlx.Document) (err error) {
	days := make(map[string]int)
	days["Monday"] = 1
	days["Tuesday"] = 2
	days["Wednesday"] = 3
	days["Thursday"] = 4
	days["Friday"] = 5
	days["Saturday"] = 6
	days["Sunday"] = 7

	getChan := func(pubdate, title string) *Channel {
		for _, c := range this.Channels {
			switch {
			case len(pubdate) > 0:
				if c.PubDate == pubdate {
					return c
				}
			case len(title) > 0:
				if c.Title == title {
					return c
				}
			}
		}
		return nil
	}

	var ch *Channel
	var i *Item
	var n *xmlx.Node
	var list, tl []*xmlx.Node
	const ns = "*"

	root := doc.SelectNode(ns, "rss")
	if root == nil {
		root = doc.SelectNode(ns, "RDF")
	}

	if root == nil {
		return errors.New("Failed to find rss/rdf node in XML.")
	}

	channels := root.SelectNodes(ns, "channel")
	for _, node := range channels {
		if ch = getChan(node.S(ns, "pubDate"), node.S(ns, "title")); ch == nil {
			ch = new(Channel)
			this.Channels = append(this.Channels, ch)
		}

		ch.Title = node.S(ns, "title")
		list = node.SelectNodes(ns, "link")
		ch.Links = make([]Link, len(list))

		for i, v := range list {
			ch.Links[i].Href = v.GetValue()
		}

		ch.Description = node.S(ns, "description")
		ch.Language = node.S(ns, "language")
		ch.Copyright = node.S(ns, "copyright")
		ch.ManagingEditor = node.S(ns, "managingEditor")
		ch.WebMaster = node.S(ns, "webMaster")
		ch.PubDate = node.S(ns, "pubDate")
		ch.LastBuildDate = node.S(ns, "lastBuildDate")
		ch.Docs = node.S(ns, "docs")

		list = node.SelectNodes(ns, "category")
		ch.Categories = make([]*Category, len(list))
		for i, v := range list {
			ch.Categories[i] = new(Category)
			ch.Categories[i].Domain = v.As(ns, "domain")
			ch.Categories[i].Text = v.GetValue()
		}

		if n = node.SelectNode(ns, "generator"); n != nil {
			ch.Generator = Generator{}
			ch.Generator.Text = n.GetValue()
		}

		ch.TTL = node.I(ns, "ttl")
		ch.Rating = node.S(ns, "rating")

		list = node.SelectNodes(ns, "hour")
		ch.SkipHours = make([]int, len(list))
		for i, v := range list {
			ch.SkipHours[i] = v.I(ns, "hour")
		}

		list = node.SelectNodes(ns, "days")
		ch.SkipDays = make([]int, len(list))
		for i, v := range list {
			ch.SkipDays[i] = days[v.GetValue()]
		}

		if n = node.SelectNode(ns, "image"); n != nil {
			ch.Image.Title = n.S(ns, "title")
			ch.Image.Url = n.S(ns, "url")
			ch.Image.Link = n.S(ns, "link")
			ch.Image.Width = n.I(ns, "width")
			ch.Image.Height = n.I(ns, "height")
			ch.Image.Description = n.S(ns, "description")
		}

		if n = node.SelectNode(ns, "cloud"); n != nil {
			ch.Cloud = Cloud{}
			ch.Cloud.Domain = n.As(ns, "domain")
			ch.Cloud.Port = n.Ai(ns, "port")
			ch.Cloud.Path = n.As(ns, "path")
			ch.Cloud.RegisterProcedure = n.As(ns, "registerProcedure")
			ch.Cloud.Protocol = n.As(ns, "protocol")
		}

		if n = node.SelectNode(ns, "textInput"); n != nil {
			ch.TextInput = Input{}
			ch.TextInput.Title = n.S(ns, "title")
			ch.TextInput.Description = n.S(ns, "description")
			ch.TextInput.Name = n.S(ns, "name")
			ch.TextInput.Link = n.S(ns, "link")
		}

		itemcount := len(ch.Items)
		list = node.SelectNodes(ns, "item")
		if len(list) == 0 {
			list = doc.SelectNodes(ns, "item")
		}

		for _, item := range list {
			i = new(Item)
			i.Title = item.S(ns, "title")
			i.Description = item.S(ns, "description")

			tl = item.SelectNodes(ns, "link")
			for _, v := range tl {
				lnk := new(Link)
				lnk.Href = v.GetValue()
				i.Links = append(i.Links, lnk)
			}

			if n = item.SelectNode(ns, "author"); n != nil {
				i.Author.Name = n.GetValue()

			} else if n = item.SelectNode(ns, "creator"); n != nil {
				i.Author.Name = n.GetValue()
			}

			i.Comments = item.S(ns, "comments")
			
			guid := item.S(ns, "guid")
			if len(guid) > 0 {
				i.Guid = &guid
			}

			i.PubDate = item.S(ns, "pubDate")

			tl = item.SelectNodes(ns, "category")
			for _, lv := range tl {
				cat := new(Category)
				cat.Domain = lv.As(ns, "domain")
				cat.Text = lv.GetValue()
				i.Categories = append(i.Categories, cat)
			}

			tl = item.SelectNodes(ns, "enclosure")
			for _, lv := range tl {
				enc := new(Enclosure)
				enc.Url = lv.As(ns, "url")
				enc.Length = lv.Ai64(ns, "length")
				enc.Type = lv.As(ns, "type")
				i.Enclosures = append(i.Enclosures, enc)
			}

			if src := item.SelectNode(ns, "source"); src != nil {
				i.Source = new(Source)
				i.Source.Url = src.As(ns, "url")
				i.Source.Text = src.GetValue()
			}

			tl = item.SelectNodes("http://purl.org/rss/1.0/modules/content/", "*")
			for _, lv := range tl {
				if lv.Name.Local == "encoded" {
					i.Content = new(Content)
					i.Content.Text = lv.String()
					break
				}
			}

			ch.Items = append(ch.Items, i)
		}

		if itemcount != len(ch.Items) && this.itemhandler != nil {
			this.itemhandler(this, ch, ch.Items[itemcount:])
		}
	}
	return
}
