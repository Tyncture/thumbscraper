# thumbscraper
A web image scraper built in Go that can extract all image URLs and/or determine 
primary images of a web page for thumbnail generation. It uses the `colly` scraper 
to scrape elements from the DOM.

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

`ImageNode` represents information relating to HTML images elements discovered on
the requested URLs.

#### func  GetImageNodes

```go
func GetImageNodes(pageURL string) ([]ImageNode, error)
```
`GetImageNodes` returns an `[]ImageNode` containing the names, alt tags, URLs and
whether an image is from an OpenGraph image meta tag.

#### type ImageNodeInfo

```go
type ImageNodeInfo struct {
	ImageNode
	Format string
	Height int
	Width  int
}
```

`ImageNodeInfo` represents information relating to image elements discovered on
the requested URLs with additional useful information.

#### func  GetImageNodeInfo

```go
func GetImageNodeInfo(imageNode ImageNode) (*ImageNodeInfo, error)
```
`GetImageNodeInfo` takes an `ImageNode` and returns an `*ImageNodeInfo` struct with
additional properties received after loading and analysing the image itself

#### func  GetImageNodeInfoBatch

```go
func GetImageNodeInfoBatch(imageNodes []ImageNode,
	requireAll ...bool) ([]*ImageNodeInfo, error)
```
`GetImageNodeInfoBatch` does the same thing as `GetImageNodeInfo`, but takes an
`ImageNode[]` instead to allow you to get an `[]ImageNodeInfo` back after processing
them in batch. The last parameter, requireAll, is an optional parameter that
will allow you to force this function to return an error if not all image nodes
could be processed. By default, it will not return an error on partial success.

#### func  EnforceURLSchema

```go
func EnforceURLSchema(pageURL string, imageURL string) string
```
EnforceURLSchema enforces the proper URL format to allow requests to be made to
retrieve them. Images embeded in HTML image elements are often missing the
schema prefix.

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
