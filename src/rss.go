package feeder

import "os"
import "xmlx"

func (this *Feed) readRss2(doc *xmlx.Document) (err os.Error) {
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

	haveItem := func(ch *Channel, pubdate, title, desc string) bool {
		for _, item := range ch.Items {
			switch {
			case len(pubdate) > 0:
				if item.PubDate == pubdate {
					return true
				}
			case len(title) > 0:
				if item.Title == title {
					return true
				}
			case len(desc) > 0:
				if item.Description == desc {
					return true
				}
			}
		}
		return false
	}

	var ch *Channel
	var i *Item
	var n *xmlx.Node
	var list, tl []*xmlx.Node

	channels := doc.SelectNodes("", "channel")
	for _, node := range channels {
		if ch = getChan(node.GetValue("", "pubDate"), node.GetValue("", "title")); ch == nil {
			ch = new(Channel)
			this.Channels = append(this.Channels, ch)
		}

		ch.Title = node.GetValue("", "title")
		list = node.SelectNodes("", "link")
		ch.Links = make([]Link, len(list))

		for i, v := range list {
			ch.Links[i].Href = v.Value
		}

		ch.Description = node.GetValue("", "description")
		ch.Language = node.GetValue("", "language")
		ch.Copyright = node.GetValue("", "copyright")
		ch.ManagingEditor = node.GetValue("", "managingEditor")
		ch.WebMaster = node.GetValue("", "webMaster")
		ch.PubDate = node.GetValue("", "pubDate")
		ch.LastBuildDate = node.GetValue("", "lastBuildDate")
		ch.Docs = node.GetValue("", "docs")

		list = node.SelectNodes("", "category")
		ch.Categories = make([]*Category, len(list))
		for i, v := range list {
			ch.Categories[i] = new(Category)
			ch.Categories[i].Domain = v.GetAttr("", "domain")
			ch.Categories[i].Text = v.Value
		}

		if n = node.SelectNode("", "generator"); n != nil {
			ch.Generator = Generator{}
			ch.Generator.Text = n.Value
		}

		ch.TTL = node.GetValuei("", "ttl")
		ch.Rating = node.GetValue("", "rating")

		list = node.SelectNodes("", "hour")
		ch.SkipHours = make([]int, len(list))
		for i, v := range list {
			ch.SkipHours[i] = int(v.GetValuei("", "hour"))
		}

		list = node.SelectNodes("", "days")
		ch.SkipDays = make([]int, len(list))
		for i, v := range list {
			ch.SkipDays[i] = days[v.Value]
		}

		if n = node.SelectNode("", "image"); n != nil {
			ch.Image.Title = n.GetValue("", "title")
			ch.Image.Url = n.GetValue("", "url")
			ch.Image.Link = n.GetValue("", "link")
			ch.Image.Width = n.GetValuei("", "width")
			ch.Image.Height = n.GetValuei("", "height")
			ch.Image.Description = n.GetValue("", "description")
		}

		if n = node.SelectNode("", "cloud"); n != nil {
			ch.Cloud = Cloud{}
			ch.Cloud.Domain = n.GetAttr("", "domain")
			ch.Cloud.Port = n.GetAttri("", "port")
			ch.Cloud.Path = n.GetAttr("", "path")
			ch.Cloud.RegisterProcedure = n.GetAttr("", "registerProcedure")
			ch.Cloud.Protocol = n.GetAttr("", "protocol")
		}

		if n = node.SelectNode("", "textInput"); n != nil {
			ch.TextInput = Input{}
			ch.TextInput.Title = n.GetValue("", "title")
			ch.TextInput.Description = n.GetValue("", "description")
			ch.TextInput.Name = n.GetValue("", "name")
			ch.TextInput.Link = n.GetValue("", "link")
		}

		itemcount := len(ch.Items)
		list = node.SelectNodes("", "item")

		for _, item := range list {
			if haveItem(ch, item.GetValue("", "pubDate"),
				item.GetValue("", "title"), item.GetValue("", "description")) {
				continue
			}

			i = new(Item)
			i.Title = item.GetValue("", "title")
			i.Description = item.GetValue("", "description")

			tl = node.SelectNodes("", "link")
			for _, v := range tl {
				lnk := new(Link)
				lnk.Href = v.Value
				i.Links = append(i.Links, lnk)
			}

			if n = item.SelectNode("", "author"); n != nil {
				i.Author = Author{}
				i.Author.Name = n.Value
			}

			i.Comments = item.GetValue("", "comments")
			i.Guid = item.GetValue("", "guid")
			i.PubDate = item.GetValue("", "pubDate")

			tl = item.SelectNodes("", "category")
			for _, lv := range tl {
				cat := new(Category)
				cat.Domain = lv.GetAttr("", "domain")
				cat.Text = lv.Value
				i.Categories = append(i.Categories, cat)
			}

			tl = item.SelectNodes("", "enclosure")
			for _, lv := range tl {
				enc := new(Enclosure)
				enc.Url = lv.GetAttr("", "url")
				enc.Length = lv.GetAttri64("", "length")
				enc.Type = lv.GetAttr("", "type")
				i.Enclosures = append(i.Enclosures, enc)
			}

			if src := item.SelectNode("", "source"); src != nil {
				i.Source = new(Source)
				i.Source.Url = src.GetAttr("", "url")
				i.Source.Text = src.Value
			}

			ch.Items = append(ch.Items, i)
		}

		if itemcount != len(ch.Items) && this.itemhandler != nil {
			this.itemhandler(this, ch, ch.Items[itemcount:])
		}
	}
	return
}
