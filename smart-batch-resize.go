package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/artyom/smartcrop"
	"github.com/nfnt/resize"
)

func main() {
	inDirPtr := flag.String("in", "./", "The input file.")
	outDirPtr := flag.String("out", "../resized", "The output file.")
	heightPtr := flag.Float64("height", 100, "The height of the image.")
	widthPtr := flag.Float64("width", 100, "The width of the image.")

	flag.Parse()

	files, _ := ioutil.ReadDir(*inDirPtr)
	r, _ := regexp.Compile(".jpe?g$")

	if _, err := os.Stat(*outDirPtr); os.IsNotExist(err) {
		fmt.Println(*outDirPtr, "folder does not exist, making")
		os.Mkdir(*outDirPtr, 0777)
	}
	for _, f := range files {
		if r.MatchString(f.Name()) {
			file, err := os.Open(*inDirPtr + "/" + f.Name())
			handleError(err)
			// decode jpeg into image.Image
			img, err := jpeg.Decode(file)
			handleError(err)
			file.Close()

			resizedImage := resizeImage(img, *heightPtr, *widthPtr)

			newImageName := r.ReplaceAllString(f.Name(), ".jpg")

			fmt.Println(f.Name(), "-->", newImageName)

			out, err := os.Create(*outDirPtr + "/" + newImageName)
			handleError(err)
			defer out.Close()

			// write new image to file
			jpeg.Encode(out, resizedImage, nil)
		} else {
			fmt.Println(f.Name(), "isn't a picture, skipping.")
		}
	}
}

// SubImager acts as a way of gathering a sub-matrix of an image.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func resizeImage(img image.Image, newHeight float64, newWidth float64) (resizedImage image.Image) {
	croppedImage := resize.Resize(uint(newWidth), 0, img, resize.Lanczos3)

	topCrop, err := smartcrop.Crop(croppedImage, int(newWidth), int(newHeight))
	handleError(err)
	fmt.Printf("Top crop: %+v\n", topCrop)

	sub, ok := croppedImage.(SubImager)
	if !ok {
		fmt.Println("No SubImage support")
		os.Exit(1)
	}

	resizedImage = sub.SubImage(topCrop)
	return
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
