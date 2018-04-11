package image

import (
	"image"
	"image/color"
	"os"
	"image/draw"
	"image/png"
	"net/http"
	"io"
	"github.com/twinj/uuid"
	"github.com/nfnt/resize"
	"encoding/base64"
	"io/ioutil"
	"../helper"
	"../model"
)

type Pixel struct {
	Point image.Point
	Color color.Color
}

func CreateFullDrawing(drawing model.Drawing) (string) {
	var urls = make([]string, 3)

	urls[0] = helper.GetStringFromNullable(drawing.Top.Url)
	urls[1] = helper.GetStringFromNullable(drawing.Middle.Url)
	urls[2] = helper.GetStringFromNullable(drawing.Bottom.Url)

	return combineImages(urls)
}

// Keep it DRY so don't have to repeat opening file and decode
func openAndDecode(url string) (image.Image, string, error) {

	response, err := http.Get(url)
	helper.CheckError(err)
	defer response.Body.Close()

	filepath := "./" + uuid.NewV4().String() + ".png"
	//open a file for writing
	file, err := os.Create(filepath)
	helper.CheckError(err)
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	helper.CheckError(err)
	file.Close()

	imgFile, err := os.Open(filepath)
	helper.CheckError(err)

	img, format, err := image.Decode(imgFile)
	helper.CheckError(err)
	imgFile.Close()

	err = os.Remove(filepath)
	helper.CheckError(err)

	modified := resize.Resize(640, 0, img, resize.Lanczos3)

	return modified, format, nil
}

// Decode image.Image's pixel data into []*Pixel
func decodePixelsFromImage(img image.Image, offsetX, offsetY int) []*Pixel {
	var pixels []*Pixel
	for y := 0; y <= img.Bounds().Max.Y; y++ {
		for x := 0; x <= img.Bounds().Max.X; x++ {
			p := &Pixel{
				Point: image.Point{X: x + offsetX, Y: y + offsetY},
				Color: img.At(x, y),
			}
			pixels = append(pixels, p)
		}
	}
	return pixels
}

//https://stackoverflow.com/questions/35964656/golang-how-to-concatenate-append-images-to-one-another
func combineImages(urls []string) string {
	var images = make([]image.Image, 0)
	var pixels = make([]*Pixel, 0)

	for _, url := range urls {
		img, _, err := openAndDecode(url)
		if err != nil {
			panic(err)
		}
		images = append(images, img)
	}

	yOffset := 0
	maxWidth := 0

	for _, img := range images {
		pixels = append(pixels, decodePixelsFromImage(img, 0, yOffset)...)

		if maxWidth == 0 || img.Bounds().Max.X > maxWidth {
			maxWidth = img.Bounds().Max.X
		}

		yOffset += img.Bounds().Max.Y
	}

	// Set a new size for the new image equal to the max width
	// of bigger image and max height of two images combined
	newRect := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: maxWidth,
			Y: yOffset,
		},
	}
	finImage := image.NewRGBA(newRect)
	// This is the cool part, all you have to do is loop through
	// each Pixel and set the image's color on the go
	for _, px := range pixels {
		finImage.Set(
			px.Point.X,
			px.Point.Y,
			px.Color,
		)
	}
	draw.Draw(finImage, finImage.Bounds(), finImage, image.Point{0, 0}, draw.Src)

	// Create a new file, write to it and then remove it

	filepath := "./" + uuid.NewV4().String() + ".png"
	out, err := os.Create(filepath)
	helper.CheckError(err)

	err = png.Encode(out, finImage)
	helper.CheckError(err)

	b, err := ioutil.ReadFile(filepath)
	helper.CheckError(err)

	out.Close()

	err = os.Remove(filepath)
	helper.CheckError(err)

	return base64.StdEncoding.EncodeToString(b)
}
