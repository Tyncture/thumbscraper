# thumbscraper
![Repository Size](https://img.shields.io/github/repo-size/Tyncture/thumbscraper.svg?t&style=flat-square)
![License](https://img.shields.io/github/license/Tyncture/thumbscraper.svg?&style=flat-square)
![Top Language](https://img.shields.io/github/languages/top/Tyncture/thumbscraper.svg?&style=flat-square)
[![GoDoc](https://godoc.org/github.com/Tyncture/thumbscraper?status.svg)](https://godoc.org/github.com/Tyncture/thumbscraper)

A web image scraper built in Go that can extract all image URLs and/or determine 
primary images of a web page for thumbnail generation. It uses the `colly` scraper 
to scrape elements from the DOM.

## Installation

Simply run the following command to install the package to your `$GOPATH`.
```go
go get "github.com/tyncture/thumbscraper"
```


And then you can use it in your project like so.
```go
package main

import (
	"github.com/tyncture/thumbscraper"
)

func main() {
	imageNodes, err := thumbscraper.GetImageNodes("https://github.com/Tyncture/thumbscraper")
	if err != nil {
		// Failed to load the web page
		// Your error handling here
	}

	imageNodesInfo, err := thumbscraper.GetImageNodeInfoBatch(imageNodes)
	if err != nil {
		// If using RequireAll is supplied in the optional GetImageNodeInfoBatchOptions
		// and is set to true, GetImageNodeInfoBatch will return an error if it cannot 
		// retrieve all images.
		// Otherwise, if it is  false or the options argument is empty, it will run
		// to completion, even if some images cannot be processed successfully
		// Your error handling here
	}
}
```

## Documentation

#### type ImageNode

```go
type ImageNode struct {
	Name           string
	Alt            string
	URL            string
	OpenGraphImage bool
}
```

ImageNode represents information relating to HTML images elements discovered on
the requested URLs.

#### func  GetImageNodes

```go
func GetImageNodes(pageURL string) ([]ImageNode, error)
```
GetImageNodes returns an []ImageNode containing the names, alt tags, URLs and
whether an image is from an OpenGraph image meta tag.

#### type ImageNodeInfo

```go
type ImageNodeInfo struct {
	ImageNode
	Format string
	Height int
	Width  int
	Image  *image.Image
}
```

ImageNodeInfo represents information relating to image elements discovered on
the requested URLs with additional useful information. Image is only populated
if ScrapeImages is set to true in ImageNodeInfoOptions or
ImageNodeInfoBatchOptions.

#### func  GetImageNodeInfo

```go
func GetImageNodeInfo(imageNode ImageNode, options ...GetImageNodeInfoOptions) (*ImageNodeInfo, error)
```
GetImageNodeInfo takes an ImageNode and returns an *ImageNodeInfo struct with
additional properties received after loading and analysing the image itself.
options is an optional GetImageNodeInfoOptions struct to specify whether to keep
images in the returned ImageNodeInfo struct, default of which is false.

#### func  GetImageNodeInfoBatch

```go
func GetImageNodeInfoBatch(imageNodes []ImageNode,
	options ...GetImageNodeInfoBatchOptions) ([]*ImageNodeInfo, error)
```
GetImageNodeInfoBatch does the same thing as GetImageNodeInfo, but takes an
ImageNode[] instead to allow you to get an []ImageNodeInfo back after processing
them in batch. options is an optional GetImageNodeInfoBatch options struct to
specify whether to keep images in the returned ImageNodeInfo structs, default of
which is false, and whether to require all image requests to complete
successfully, default of which is also false. Refer to struct type
GetImageNodeInfoBatchOptions for more information.

#### type GetImageNodeInfoBatchOptions

```go
type GetImageNodeInfoBatchOptions struct {
	GetImageNodeInfoOptions
	RequireAll bool
}
```

GetImageNodeInfoBatchOptions represents the configuration used by
GetImageNodeInfoBatch. Default for RequireAll is false.

#### type GetImageNodeInfoOptions

```go
type GetImageNodeInfoOptions struct {
	ScrapeImages bool
}
```

GetImageNodeInfoOptions represents the configuration used by GetImageNodeInfo.
Default for ScrapeImages is false.

#### func  EnforceURLSchema

```go
func EnforceURLSchema(pageURL string, imageURL string) string
```
EnforceURLSchema enforces the proper URL format to allow requests to be made to
retrieve them. Images embeded in HTML image elements are often missing the
schema prefix. This is used by GetImageNodeInfo to ensure that the URL is valid 
before making a request for the image resource.

## License
```
MIT License

Copyright (c) 2019 John Su

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
