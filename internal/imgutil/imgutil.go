package imgutil

import (
	"bytes"
	"image"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

// Load will load an image from file
func Load(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

// SavePNG will save an image to a PNG file
func SavePNG(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}

// SaveGIF returns an animated GIF given a series of frames
func SaveGIF(path string, frames []image.Image, delay int) {
	outGIF := gif.GIF{
		Image: make([]*image.Paletted, len(frames)),
		Delay: make([]int, len(frames)),
	}

	// convert each frame to a palleted image within the GIF
	for i := range frames {
		buf := bytes.Buffer{}
		gif.Encode(&buf, frames[i], nil)
		palettedFrame, err := gif.Decode(&buf)
		if err != nil {
			log.Fatal(err)
		}

		outGIF.Image[i] = palettedFrame.(*image.Paletted)
		outGIF.Delay[i] = delay
	}

	// open output file
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// write GIF to file
	err = gif.EncodeAll(f, &outGIF)
	if err != nil {
		log.Fatal(err)
	}
}
