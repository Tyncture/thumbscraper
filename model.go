package thumbscraper

import (
	"image"
)

// ImageNode represents information relating to
// HTML images elements discovered on the requested URLs.
type ImageNode struct {
	Name           string
	Alt            string
	URL            string
	OpenGraphImage bool
}

// ImageNodeInfo represents information relating to
// image elements discovered on the requested URLs with
// additional useful information. Image is only populated
// if ScrapeImages is set to true in ImageNodeInfoOptions
// or ImageNodeInfoBatchOptions.
type ImageNodeInfo struct {
	ImageNode
	Format string
	Height int
	Width  int
	Image  *image.Image
}
