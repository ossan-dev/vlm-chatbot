package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"

	"golang.org/x/image/draw"
)

const MAX_NUM_PIXELS = 768

func ScaleDownImage(img image.Image) image.Image {
	// resize the image if needed
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// check if we exceed one of the boundaries
	if width > MAX_NUM_PIXELS || height > MAX_NUM_PIXELS {
		var newWidth, newHeight int
		// landscape scenario
		if width > height {
			newWidth = MAX_NUM_PIXELS
			newHeight = height * MAX_NUM_PIXELS / width
		} else {
			//portrait scenario
			newHeight = MAX_NUM_PIXELS
			newWidth = width * MAX_NUM_PIXELS / height
		}

		resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.CatmullRom.Scale(resized, resized.Rect, img, img.Bounds(), draw.Over, nil)
		img = resized
		fmt.Printf("🖼️ Image resized from %dx%d to %dx%d\n", bounds.Dx(), bounds.Dy(), newWidth, newHeight)
	}
	return img
}

func main() {
	// open the image
	file, err := os.Open("/home/ossan/Projects/pdf-chatbot/imgs/booking.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	scaledImg := ScaleDownImage(img)

	// encode to PNG
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	var buf bytes.Buffer
	err = encoder.Encode(&buf, scaledImg)
	if err != nil {
		panic(err)
	}
	newImg, err := os.Create("booking_resized.png")
	if err != nil {
		panic(err)
	}
	defer newImg.Close()

	writtenChars, err := io.Copy(newImg, &buf)
	if err != nil {
		panic(err)
	}
	fmt.Println("Written Characters:", writtenChars)

}
