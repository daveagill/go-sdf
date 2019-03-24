package main

import (
	"flag"
	"image"
	"image/color"
	"log"

	"github.com/daveagill/go-sdf/internal/imgutil"
	"github.com/daveagill/go-sdf/sdf"
)

func main() {
	var (
		startPath  string
		endPath    string
		outPath    string
		numFrames  int
		frameDelay int
		boomerang  bool
	)

	flag.StringVar(&startPath, "from", "", "Required. The start image (png/jpeg/gif)")
	flag.StringVar(&endPath, "to", "", "Required. The end image (png/jpeg/gif)")
	flag.StringVar(&outPath, "out", "out.gif", "The output gif")
	flag.IntVar(&numFrames, "frames", 10, "The number of frames to generate")
	flag.IntVar(&frameDelay, "framedelay", 0, "The delay between frames in 100ths of a second")
	flag.BoolVar(&boomerang, "boomerang", true, "Whether to animate back to the initial image (Doubles the number of frames)")
	flag.Parse()

	if startPath == "" {
		log.Fatal("-from not specified")
	}
	if endPath == "" {
		log.Fatal("-to not specified")
	}

	startImg := imgutil.Load(startPath)
	endImg := imgutil.Load(endPath)

	startSDF := sdf.FromImageAlpha(startImg, sdf.HalfAlpha)
	endSDF := sdf.FromImageAlpha(endImg, sdf.HalfAlpha)

	if startSDF.Width != endSDF.Width || startSDF.Height != endSDF.Height {
		log.Fatal("Images do not have the same dimensions")
	}

	// animate from start to end
	frames := make([]image.Image, numFrames, numFrames*2)
	for i := 0; i < numFrames; i++ {
		field, _ := sdf.Lerp(startSDF, endSDF, float64(i)/float64(numFrames-1))
		frames[i] = field.DrawImplicitSurface(0, color.Black, color.White)
	}

	// create the reverse sequence of frames to 'boomerang' back to the start
	if boomerang {
		for i := numFrames - 1; i >= 0; i-- {
			frames = append(frames, frames[i])
		}
	}

	imgutil.SaveGIF(outPath, frames, frameDelay)
}
