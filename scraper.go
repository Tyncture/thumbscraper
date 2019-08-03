package thumbscraper

import (
	"errors"
	"fmt"
	"image"

	// Image type import as side effects to support
	// different image formats
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var httpSchemaRegex *regexp.Regexp
var urlStartRegex *regexp.Regexp

// ImageNodeInfoOptions represents the configuration used by
// GetImageNodeInfo. Default for ScrapeImages is false.
type ImageNodeInfoOptions struct {
	ScrapeImages bool
}

// ImageNodeInfoBatchOptions represents the configuration used by
// GetImageNodeInfoBatch. Default for RequireAll is false.
type ImageNodeInfoBatchOptions struct {
	ImageNodeInfoOptions
	RequireAll bool
}

// GetImageNodes returns an []ImageNode containing the
// names, alt tags, URLs and whether an image is from an
// OpenGraph image meta tag.
func GetImageNodes(pageURL string) ([]ImageNode, error) {
	var err error
	imageNodes := []ImageNode{}

	c := colly.NewCollector()

	// Get images specified in the OpenGraph og:image meta property tags,
	// which are used by sites such as Facebook and Reddit to determine
	// which image from a web page should be featured as the main thumbnail
	c.OnHTML("meta[property=\"og:image\"][content]", func(e *colly.HTMLElement) {
		imgURL := EnforceURLSchema(pageURL, e.Attr("content"))
		imgURLSegments := filterEmptyStrings(strings.Split(imgURL, "/"))
		name := imgURLSegments[len(imgURLSegments)-1]

		if len(pageURL) > 0 {
			imageNode := ImageNode{
				Name:           name,
				Alt:            "",
				URL:            imgURL,
				OpenGraphImage: true,
			}
			imageNodes = append(imageNodes, imageNode)
		}
	})

	// Get all other images which simply use a HTML image element and
	// scrape their useful properties
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		alt := e.Attr("alt")
		imgURL := EnforceURLSchema(pageURL, e.Attr("src"))
		imgURLSegments := filterEmptyStrings(strings.Split(imgURL, "/"))
		name := imgURLSegments[len(imgURLSegments)-1]

		if len(pageURL) > 0 {
			imageNode := ImageNode{
				Name:           name,
				Alt:            alt,
				URL:            imgURL,
				OpenGraphImage: false,
			}
			imageNodes = append(imageNodes, imageNode)
		}
	})

	// Return the error from colly if it manifests
	c.OnError(func(_ *colly.Response, cErr error) {
		err = cErr
	})
	c.Visit(pageURL)

	return imageNodes, err
}

// GetImageNodeInfo takes an ImageNode and returns an *ImageNodeInfo
// struct with additional properties received after loading and
// analysing the image itself
func GetImageNodeInfo(imageNode ImageNode) (*ImageNodeInfo, error) {
	res, err := http.Get(imageNode.URL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 || res.StatusCode == 201 {
		return nil, errors.New(res.Status)
	}

	img, format, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	imageBounds := img.Bounds()
	imageNodeInfo := &ImageNodeInfo{
		ImageNode: imageNode,
		Format:    format,
		Width:     imageBounds.Max.X,
		Height:    imageBounds.Max.Y,
	}
	return imageNodeInfo, nil
}

// GetImageNodeInfoBatch does the same thing as GetImageNodeInfo,
// but takes an ImageNode[] instead to allow you to get an
// []ImageNodeInfo back after processing them in batch. The last
// parameter, requireAll, is an optional parameter that will allow
// you to force this function to return an error if not all image
// nodes could be processed. By default, it will not return an error
// on partial success.
func GetImageNodeInfoBatch(imageNodes []ImageNode,
	requireAll ...bool) ([]*ImageNodeInfo, error) {
	imageNodesInfo := []*ImageNodeInfo{}

	for _, imageNode := range imageNodes {
		imageNodeInfo, err := GetImageNodeInfo(imageNode)
		if err == nil {
			imageNodesInfo = append(imageNodesInfo, imageNodeInfo)
		} else if len(requireAll) > 0 && requireAll[0] {
			return nil, err
		}
	}
	return imageNodesInfo, nil
}

// EnforceURLSchema enforces the proper URL format to allow
// requests to be made to retrieve them. Images embeded in HTML
// image elements are often missing the schema prefix.
func EnforceURLSchema(pageURL string, imageURL string) string {
	if httpSchemaRegex.MatchString(imageURL) {
		return imageURL
	}

	urlStart := urlStartRegex.FindString(pageURL)
	urlEnd := imageURL
	if strings.HasPrefix(imageURL, "//") {
		urlStart = "https:"
	}
	if !strings.HasPrefix(imageURL, "/") {
		urlEnd = fmt.Sprintf("/%s", imageURL)
	}

	finalURL := fmt.Sprintf("%s%s", urlStart, urlEnd)
	return finalURL
}

// strings.Split returns everything before and after the delimiter,
// and in cases where there is nothing before or after it,
// it will return an empty string.
func filterEmptyStrings(strings []string) []string {
	final := []string{}
	for _, s := range strings {
		if s != "" {
			final = append(final, s)
		}
	}
	return final
}

func init() {
	httpSchemaRegex = regexp.MustCompile(`^https{0,1}:\/\/`)
	urlStartRegex = regexp.MustCompile(
		`^https{0,1}:\/\/[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+`)
}
