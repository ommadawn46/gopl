package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
)

var outformatOpt = flag.String("f", "jpg", "output image format")

func main() {
	flag.Parse()

	img, err := getImage(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(*outformatOpt) {
	case "jpg", "jpeg":
		err = outJPEG(os.Stdout, img)
	case "png":
		err = outPNG(os.Stdout, img)
	case "gif":
		err = outGIF(os.Stdout, img)
	default:
		err = fmt.Errorf("unknown image format")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getImage(in io.Reader) (image.Image, error) {
	img, kind, err := image.Decode(in)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	return img, nil
}

func outJPEG(out io.Writer, img image.Image) error {
	if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 95}); err != nil {
		return fmt.Errorf("jpeg: %v", err)
	}
	return nil
}

func outPNG(out io.Writer, img image.Image) error {
	if err := png.Encode(os.Stdout, img); err != nil {
		return fmt.Errorf("png: %v", err)
	}
	return nil
}

func outGIF(out io.Writer, img image.Image) error {
	if err := gif.Encode(os.Stdout, img, &gif.Options{}); err != nil {
		return fmt.Errorf("gif: %v", err)
	}
	return nil
}
