package feeder

import "os"
import "xmlx"

func (this *Feed) readRss2(doc *xmlx.Document) (err os.Error) {
	channels := doc.SelectNodes("", "channel");
	for _, node := range channels {
		ch := Channel{};
		ch.Title = node.GetValue("", "title");

		list := node.SelectNodes("", "link");
		ch.Links = make([]Link, len(list));
		for i, v := range list {
			ch.Links[i].Href = v.Value;
		}

		ch.Description = node.GetValue("", "description");
		ch.Language = node.GetValue("", "language");
		ch.Copyright = node.GetValue("", "copyright");
		ch.ManagingEditor = node.GetValue("", "managingEditor");
		ch.WebMaster = node.GetValue("", "webMaster");
		ch.PubDate = node.GetValue("", "pubDate");
		ch.LastBuildDate = node.GetValue("", "lastBuildDate");
		ch.Docs = node.GetValue("", "docs");

		list = node.SelectNodes("", "category");
		ch.Categories = make([]Category, len(list));
		for i, v := range list {
			ch.Categories[i].Domain = v.GetAttr("", "domain");
			ch.Categories[i].Text = v.Value;
		}

		n := node.SelectNode("", "generator");
		if n != nil {
			ch.Generator = Generator{};
			ch.Generator.Text = n.Value;
		}

		ch.TTL = node.GetValuei("", "ttl");
		ch.Rating = node.GetValue("", "rating");

		list = node.SelectNodes("", "hour");
		ch.SkipHours = make([]int, len(list));
		for i, v := range list {
			ch.SkipHours[i] = int(v.GetValuei("", "hour"));
		}

		list = node.SelectNodes("", "days");
		ch.SkipDays = make([]int, len(list));
		for i, v := range list {
			ch.SkipDays[i] = mapDay(v.Value);
		}

		n = node.SelectNode("", "image");
		if n != nil {
			ch.Image.Title = n.GetValue("", "title");
			ch.Image.Url = n.GetValue("", "url");
			ch.Image.Link = n.GetValue("", "link");
			ch.Image.Width = n.GetValuei("", "width");
			ch.Image.Height = n.GetValuei("", "height");
			ch.Image.Description = n.GetValue("", "description");
		}

		n = node.SelectNode("", "cloud");
		if n != nil {
			ch.Cloud = Cloud{};
			ch.Cloud.Domain = n.GetAttr("", "domain");
			ch.Cloud.Port = n.GetAttri("", "port");
			ch.Cloud.Path = n.GetAttr("", "path");
			ch.Cloud.RegisterProcedure = n.GetAttr("", "registerProcedure");
			ch.Cloud.Protocol = n.GetAttr("", "protocol");
		}

		n = node.SelectNode("", "textInput");
		if n != nil {
			ch.TextInput = Input{};
			ch.TextInput.Title = n.GetValue("", "title");
			ch.TextInput.Description = n.GetValue("", "description");
			ch.TextInput.Name = n.GetValue("", "name");
			ch.TextInput.Link = n.GetValue("", "link");
		}

		list = node.SelectNodes("", "item");
		for _, item := range list {
			i := Item{};
			i.Title = item.GetValue("", "title");
			i.Description = item.GetValue("", "description");

			list = node.SelectNodes("", "link");
			i.Links = make([]Link, 0);
			for _, v := range list {
				lnk := Link{};
				lnk.Href = v.Value;
				i.addLink(lnk);
			}

			n = item.SelectNode("", "author");
			if n != nil {
				i.Author = Author{};
				i.Author.Name = n.Value;
			}

			i.Comments = item.GetValue("", "comments");
			i.Guid = item.GetValue("", "guid");
			i.PubDate = item.GetValue("", "pubDate");

			list := item.SelectNodes("", "category");
			i.Categories = make([]Category, len(list));
			for li, lv := range list {
				i.Categories[li].Domain = lv.GetAttr("", "domain");
				i.Categories[li].Text = lv.Value;
			}

			list = item.SelectNodes("", "enclosure");
			i.Enclosures = make([]Enclosure, len(list));
			for li, lv := range list {
				i.Enclosures[li].Url = lv.GetAttr("", "url");
				i.Enclosures[li].Length = lv.GetAttri64("", "length");
				i.Enclosures[li].Type = lv.GetAttr("", "type");
			}

			src := item.SelectNode("", "source");
			if src != nil {
				i.Source = Source{};
				i.Source.Url = src.GetAttr("", "url");
				i.Source.Text = src.Value;
			}

			ch.addItem(i);
		}

		this.addChannel(ch);
	}
	return
}

func mapDay(day string) int {
	switch day {
	case "Monday": return 1;
	case "Tuesday": return 2;
	case "Wednesday": return 3;
	case "Thursday": return 4;
	case "Friday": return 5;
	case "Saturday": return 6;
	case "Sunday": return 7;
	}
	return 1;
}

