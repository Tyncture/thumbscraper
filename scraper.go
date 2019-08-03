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
