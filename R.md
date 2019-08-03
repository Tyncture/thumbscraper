# thumbscraper
--
    import "github.com/tyncture/thumbscraper"


## Usage

#### func  EnforceURLSchema

```go
func EnforceURLSchema(pageURL string, imageURL string) string
```
EnforceURLSchema enforces the proper URL format to allow requests to be made to
retrieve them. Images embeded in HTML image elements are often missing the
schema prefix.

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
