package feeder

type Channel struct {
	Title          string
	Links          []Link
	Description    string
	Language       string
	Copyright      string
	ManagingEditor string
	WebMaster      string
	PubDate        string
	LastBuildDate  string
	Docs           string
	Categories     []Category
	Generator      Generator
	TTL            int
	Rating         string
	SkipHours      []int
	SkipDays       []int
	Image          Image
	Items          []Item
	Cloud          Cloud
	TextInput      Input

	// Atom fields
	Id       string
	Rights   string
	Author   Author
	SubTitle SubTitle
}

func (this *Channel) addItem(item Item) {
	c := make([]Item, len(this.Items)+1)
	copy(c, this.Items)
	c[len(c)-1] = item
	this.Items = c
}

func (this *Channel) addLink(l Link) {
	c := make([]Link, len(this.Links)+1)
	copy(c, this.Links)
	c[len(c)-1] = l
	this.Links = c
}
