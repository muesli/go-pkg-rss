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
	c := make([]Enclosure, len(this.Enclosures)+1)
	copy(c, this.Enclosures)
	c[len(c)-1] = e
	this.Enclosures = c
}

func (this *Item) addLink(l Link) {
	c := make([]Link, len(this.Links)+1)
	copy(c, this.Links)
	c[len(c)-1] = l
	this.Links = c
}
