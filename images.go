package rssreader

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/url"
)

type imagedata struct {
	URL           *url.URL
	Width, Height int
}

func selectBestImage(width, height int, images []*url.URL) *url.URL {
	// Check if the slice is nil or there are no images
	if images == nil || len(images) == 0 {
		return nil
	}

	// Channel for images
	imageSizes := make(chan imagedata)

	// Iterate over each image
	for _, v := range images {
		go func(imageURL *url.URL) {
			width, height := getImageDimensions(imageURL)
			imageSizes <- imagedata{
				URL:    imageURL,
				Width:  width,
				Height: height,
			}
		}(v)
	}

	// Holder
	var chosenImage *url.URL

	// Retrieve all the values one-by-one
	for i := 0; i < len(images); i++ {
		data := <-imageSizes

		// Check if the image has the required attributes
		if data.Width >= width && data.Height >= height {
			chosenImage = data.URL
			break
		}
	}

	// Close the channel
	close(imageSizes)

	return chosenImage
}

func getImageDimensions(imageURL *url.URL) (int, int) {
	// Check if image is not nil
	if imageURL == nil {
		return 0, 0
	}

	// Fetch the image
	resp, err := http.Get(imageURL.String())

	// Check if there was an error
	if err != nil {
		return 0, 0
	}

	// Close the body once we go back
	defer resp.Body.Close()

	// Try decoding the image
	imageb, _, err := image.DecodeConfig(resp.Body)

	// Check if there was an error
	if err != nil {
		return 0, 0
	}

	return imageb.Width, imageb.Height
}
