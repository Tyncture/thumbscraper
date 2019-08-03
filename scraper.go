package main

import (
	"strings"

	"github.com/gocolly/colly"
)

// GetImageNodes returns an []ImageNode containing the
// names, alt tags, URLs and whether an image is from an
// OpenGraph image meta tag
func GetImageNodes(url string) []ImageNode {
	imageNodes := []ImageNode{}

	c := colly.NewCollector()

	// Get images specified in the OpenGraph og:image meta property tags,
	// which are used by sites such as Facebook and Reddit to determine
	// which image from a web page should be featured as the main thumbnail
	c.OnHTML("meta[property=\"og:image\"][content]", func(e *colly.HTMLElement) {
		url := e.Attr("content")
		urlSegments := strings.Split(url, "/")
		name := urlSegments[len(urlSegments)-1]

		imageNode := ImageNode{
			Name:           name,
			Alt:            "",
			URL:            url,
			OpenGraphImage: true,
		}
		imageNodes = append(imageNodes, imageNode)
	})

	// Get all other images which simply use a HTML image element and
	// scrape their useful properties
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		alt := e.Attr("alt")
		url := e.Attr("src")
		urlSegments := strings.Split(url, "/")
		name := urlSegments[len(urlSegments)-1]

		imageNode := ImageNode{
			Name:           name,
			Alt:            alt,
			URL:            url,
			OpenGraphImage: false,
		}
		imageNodes = append(imageNodes, imageNode)
	})
	c.Visit(url)

	return imageNodes
}
