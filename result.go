package main

// ImageNode represents information relating to
// images elements discovered on the requested URLs
type ImageNode struct {
	Name           string
	URL            string
	OpenGraphImage bool
}

// ImageNodeInfo represents information relating to
// image elements discovered on the requested URLs with
// additional useful information
type ImageNodeInfo struct {
	ImageNode
	Height int
	Width  int
}
