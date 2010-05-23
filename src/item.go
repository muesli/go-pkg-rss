package feeder

type Item struct {
	// RSS and Shared fields
	Title       string
	Links       []Link
	Description string
	Author      Author
	Categories  []Category
	Comments    string
	Enclosures  []Enclosure
	Guid        string
	PubDate     string
	Source      Source

	// Atom specific fields
	Id           string
	Generator    Generator
	Contributors []string
	Content      Content
}

func (this *Item) addEnclosure(e Enclosure) {
	slice := make([]Enclosure, len(this.Enclosures)+1)
	for i, v := range this.Enclosures {
		slice[i] = v
	}
	slice[len(slice)-1] = e
	this.Enclosures = slice
}

func (this *Item) addLink(l Link) {
	slice := make([]Link, len(this.Links)+1)
	for i, v := range this.Links {
		slice[i] = v
	}
	slice[len(slice)-1] = l
	this.Links = slice
}
