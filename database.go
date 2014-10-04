/*
Credits go to github.com/SlyMarbo/rss for inspiring this solution.
*/
package feeder

import (
	"fmt"
)

type database struct {
	request  chan string
	response chan bool
	known    map[string]struct{}
}

func (d *database) Run() {
	d.known = make(map[string]struct{})
	var s string

	for {
		s = <-d.request
		if _, ok := d.known[s]; ok {
			fmt.Println("Database used: true")
			d.response <- true
		} else {
			fmt.Println("Database used: false")
			d.response <- false
			d.known[s] = struct{}{}
		}
	}
}

func NewDatabase() *database {
	database := new(database)
	database.request = make(chan string)
	database.response = make(chan bool)
	go database.Run()
	return database
}
