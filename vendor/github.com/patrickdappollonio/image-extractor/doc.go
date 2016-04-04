/*
Package `extractor` provides an easy option to get Images from a string of text. Could
be an HTML string or even a raw string of text. By using regular expressions, the package
iterates over the content to get any number of images that could be found.

`extractor` only works finding JPG, GIF and PNG images for now.

Example:

	imageUrl := extractor.ImageExtractor{Content: text}.GetFirst()
	fmt.Println(imageUrl)

You could also re-use the extractor, like this:

	ext := extractor.ImageExtractor{Content: text}

	imageUrl := ext.GetFirst()
	fmt.Println(imageUrl)

	imageSlice := ext.GetAll()
	fmt.Println(imageSlice)

	imageSlice = ext.GetNumber(5)
	fmt.Println(imageSlice)
*/
package extractor
