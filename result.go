package main

// ImageNode represents an image discovered on the
// requested URL during the scraping process
type ImageNode struct {
	Height         int
	Width          int
	Format         string
	URL            string
	OpenGraphImage bool
}
