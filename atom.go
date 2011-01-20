package feeder

import "os"
import xmlx "github.com/jteeuwen/go-pkg-xmlx"

func (this *Feed) readAtom(doc *xmlx.Document) (err os.Error) {
	ns := "http://www.w3.org/2005/Atom"
	channels := doc.SelectNodes(ns, "feed")

	getChan := func(id, title string) *Channel {
		for _, c := range this.Channels {
			switch {
			case len(id) > 0:
				if c.Id == id {
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

	haveItem := func(ch *Channel, id, title, desc string) bool {
		for _, item := range ch.Items {
			switch {
			case len(id) > 0:
				if item.Id == id {
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
	var tn *xmlx.Node
	var list []*xmlx.Node

	for _, node := range channels {
		if ch = getChan(node.GetValue(ns, "id"), node.GetValue(ns, "title")); ch == nil {
			ch = new(Channel)
			this.Channels = append(this.Channels, ch)
		}

		ch.Title = node.GetValue(ns, "title")
		ch.LastBuildDate = node.GetValue(ns, "updated")
		ch.Id = node.GetValue(ns, "id")
		ch.Rights = node.GetValue(ns, "rights")

		list = node.SelectNodes(ns, "link")
		ch.Links = make([]Link, len(list))
		for i, v := range list {
			ch.Links[i].Href = v.GetAttr("", "href")
			ch.Links[i].Rel = v.GetAttr("", "rel")
			ch.Links[i].Type = v.GetAttr("", "type")
			ch.Links[i].HrefLang = v.GetAttr("", "hreflang")
		}

		if tn = node.SelectNode(ns, "subtitle"); tn != nil {
			ch.SubTitle = SubTitle{}
			ch.SubTitle.Type = tn.GetAttr("", "type")
			ch.SubTitle.Text = tn.Value
		}

		if tn = node.SelectNode(ns, "generator"); tn != nil {
			ch.Generator = Generator{}
			ch.Generator.Uri = tn.GetAttr("", "uri")
			ch.Generator.Version = tn.GetAttr("", "version")
			ch.Generator.Text = tn.Value
		}

		if tn = node.SelectNode(ns, "author"); tn != nil {
			ch.Author = Author{}
			ch.Author.Name = tn.GetValue("", "name")
			ch.Author.Uri = tn.GetValue("", "uri")
			ch.Author.Email = tn.GetValue("", "email")
		}

		itemcount := len(ch.Items)
		list = node.SelectNodes(ns, "entry")

		for _, item := range list {
			if haveItem(ch, item.GetValue(ns, "id"), item.GetValue(ns, "title"), item.GetValue(ns, "summary")) {
				continue
			}

			i = new(Item)
			i.Title = item.GetValue(ns, "title")
			i.Id = item.GetValue(ns, "id")
			i.PubDate = item.GetValue(ns, "updated")
			i.Description = item.GetValue(ns, "summary")

			links := item.SelectNodes(ns, "link")
			for _, lv := range links {
				if tn.GetAttr(ns, "rel") == "enclosure" {
					enc := new(Enclosure)
					enc.Url = lv.GetAttr("", "href")
					enc.Type = lv.GetAttr("", "type")
					i.Enclosures = append(i.Enclosures, enc)
				} else {
					lnk := new(Link)
					lnk.Href = lv.GetAttr("", "href")
					lnk.Rel = lv.GetAttr("", "rel")
					lnk.Type = lv.GetAttr("", "type")
					lnk.HrefLang = lv.GetAttr("", "hreflang")
					i.Links = append(i.Links, lnk)
				}
			}

			list = item.SelectNodes(ns, "contributor")
			for _, cv := range list {
				i.Contributors = append(i.Contributors, cv.GetValue("", "name"))
			}

			if tn = item.SelectNode(ns, "content"); tn != nil {
				i.Content = new(Content)
				i.Content.Type = tn.GetAttr("", "type")
				i.Content.Lang = tn.GetValue("xml", "lang")
				i.Content.Base = tn.GetValue("xml", "base")
				i.Content.Text = tn.Value
			}

			ch.Items = append(ch.Items, i)
		}

		if itemcount != len(ch.Items) && this.itemhandler != nil {
			this.itemhandler(this, ch, ch.Items[itemcount:])
		}
	}
	return
}
