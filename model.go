package thumbscraper

// ImageNode represents information relating to
// images elements discovered on the requested URLs
type ImageNode struct {
	Name           string
	Alt            string
	URL            string
	OpenGraphImage bool
}

// ImageNodeInfo represents information relating to
// image elements discovered on the requested URLs with
// additional useful information
type ImageNodeInfo struct {
	ImageNode
	Format string
	Height int
	Width  int
}
