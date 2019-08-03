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

// GetImageNodeInfoOptions represents the configuration used by
// GetImageNodeInfo. Default for ScrapeImages is false.
type GetImageNodeInfoOptions struct {
	ScrapeImages bool
}

// GetImageNodeInfoBatchOptions represents the configuration used by
// GetImageNodeInfoBatch. Default for RequireAll is false.
type GetImageNodeInfoBatchOptions struct {
	GetImageNodeInfoOptions
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

// DetermineThumbnail returns the *ImageNodeInfo for the best thumbnail from a
// []*ImageNodeInfo. The *ImageNodeInfo will also have the image itself in the
// Image property if ScrapeImages is set to true in GetImageNodeInfoBatchOptions
// that was passed into GetImageNodeInfoBatch. error is returned if the supplied
// []*ImageNodeInfo is empty.
func DetermineThumbnail(imageNodesWithInfo []*ImageNodeInfo) (*ImageNodeInfo, error) {
	var bestMatch *ImageNodeInfo

	for _, cur := range imageNodesWithInfo {
		if cur == nil {
			continue
		}

		if cur.OpenGraphImage {
			return cur, nil
		}

		if bestMatch == nil {
			bestMatch = cur
			continue
		}

		area := cur.Width * cur.Height
		bestMatchArea := bestMatch.Width * bestMatch.Height
		if area > bestMatchArea {
			bestMatch = cur
		}
	}

	if bestMatch == nil {
		return nil, errors.New("No *ImageNodeInfo elements in supplied slice")
	}

	return bestMatch, nil
}

// GetImageNodeInfo takes an ImageNode and returns an *ImageNodeInfo
// struct with additional properties received after loading and
// analysing the image itself. options is an optional GetImageNodeInfoOptions
// struct to specify whether to keep images in the returned ImageNodeInfo
// struct, default of which is false.
func GetImageNodeInfo(imageNode ImageNode, options ...GetImageNodeInfoOptions) (*ImageNodeInfo, error) {
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
		Image:     nil,
	}

	if len(options) > 0 && options[0].ScrapeImages {
		imageNodeInfo.Image = &img
	}

	return imageNodeInfo, nil
}

// GetImageNodeInfoBatch does the same thing as GetImageNodeInfo,
// but takes an ImageNode[] instead to allow you to get an
// []ImageNodeInfo back after processing them in batch. options is
// an optional GetImageNodeInfoBatch options struct to specify whether
// to keep images in the returned ImageNodeInfo structs, default of
// which is false, and whether to require all image requests to complete
// successfully, default of which is also false. Refer to struct type
// GetImageNodeInfoBatchOptions for more information.
func GetImageNodeInfoBatch(imageNodes []ImageNode,
	options ...GetImageNodeInfoBatchOptions) ([]*ImageNodeInfo, error) {
	imageNodesInfo := []*ImageNodeInfo{}
	localOptions := GetImageNodeInfoBatchOptions{
		GetImageNodeInfoOptions: GetImageNodeInfoOptions{
			ScrapeImages: false,
		},
		RequireAll: false,
	}

	if len(options) > 0 {
		localOptions = options[0]
	}

	for _, imageNode := range imageNodes {
		imageNodeInfo, err := GetImageNodeInfo(imageNode,
			localOptions.GetImageNodeInfoOptions)
		if err == nil {
			imageNodesInfo = append(imageNodesInfo, imageNodeInfo)
		} else if localOptions.RequireAll {
			return nil, err
		}
	}
	return imageNodesInfo, nil
}

// EnforceURLSchema enforces the proper URL format to allow
// requests to be made to retrieve them. Images embeded in HTML
// image elements are often missing the schema prefix. This is used
// by GetImageNodeInfo to ensure that the URL is valid before making
// a request for the image resource.
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
