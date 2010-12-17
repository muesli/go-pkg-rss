package feeder

import "os"
import "xmlx"

func (this *Feed) readAtom(doc *xmlx.Document) (err os.Error) {
	ns := "http://www.w3.org/2005/Atom"
	channels := doc.SelectNodes(ns, "feed")
	for _, node := range channels {
		ch := Channel{}
		ch.Title = node.GetValue(ns, "title")
		ch.LastBuildDate = node.GetValue(ns, "updated")
		ch.Id = node.GetValue(ns, "id")
		ch.Rights = node.GetValue(ns, "rights")

		list := node.SelectNodes(ns, "link")
		ch.Links = make([]Link, len(list))
		for i, v := range list {
			ch.Links[i].Href = v.GetAttr("", "href")
			ch.Links[i].Rel = v.GetAttr("", "rel")
			ch.Links[i].Type = v.GetAttr("", "type")
			ch.Links[i].HrefLang = v.GetAttr("", "hreflang")
		}

		tn := node.SelectNode(ns, "subtitle")
		if tn != nil {
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

		list = node.SelectNodes(ns, "entry")
		ch.Items = make([]Item, len(list))
		for _, v := range list {
			item := Item{}
			item.Title = v.GetValue(ns, "title")
			item.Id = v.GetValue(ns, "id")
			item.PubDate = v.GetValue(ns, "updated")
			item.Description = v.GetValue(ns, "summary")

			list = v.SelectNodes(ns, "link")
			item.Links = make([]Link, 0)
			for _, lv := range list {
				if tn.GetAttr(ns, "rel") == "enclosure" {
					enc := Enclosure{}
					enc.Url = lv.GetAttr("", "href")
					enc.Type = lv.GetAttr("", "type")
					item.Enclosures = append(item.Enclosures, enc)
				} else {
					lnk := Link{}
					lnk.Href = lv.GetAttr("", "href")
					lnk.Rel = lv.GetAttr("", "rel")
					lnk.Type = lv.GetAttr("", "type")
					lnk.HrefLang = lv.GetAttr("", "hreflang")
					item.Links = append(item.Links, lnk)
				}
			}

			list = v.SelectNodes(ns, "contributor")
			item.Contributors = make([]string, len(list))
			for ci, cv := range list {
				item.Contributors[ci] = cv.GetValue("", "name")
			}

			if tn = v.SelectNode(ns, "content"); tn != nil {
				item.Content = Content{}
				item.Content.Type = tn.GetAttr("", "type")
				item.Content.Lang = tn.GetValue("xml", "lang")
				item.Content.Base = tn.GetValue("xml", "base")
				item.Content.Text = tn.Value
			}
			ch.Items = append(ch.Items, item)
		}

		this.Channels = append(this.Channels, ch)
	}
	return
}
