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
	slice := make([]Item, len(this.Items)+1)
	for i, v := range this.Items {
		slice[i] = v
	}
	slice[len(slice)-1] = item
	this.Items = slice
}


func (this *Channel) addLink(l Link) {
	slice := make([]Link, len(this.Links)+1)
	for i, v := range this.Links {
		slice[i] = v
	}
	slice[len(slice)-1] = l
	this.Links = slice
}
