package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/artyom/smartcrop"
	"github.com/nfnt/resize"
)

func main() {
	files, _ := ioutil.ReadDir("./")
	r, _ := regexp.Compile(".jpe?g$")
	if _, err := os.Stat("../resized"); os.IsNotExist(err) {
		fmt.Println("resized folder does not exist, making")
		os.Mkdir("../resized", 0777)
	}
	for _, f := range files {
		if r.MatchString(f.Name()) {
			height, err := strconv.ParseFloat(os.Args[1], 64)
			handleError(err)
			width, err := strconv.ParseFloat(os.Args[2], 64)
			handleError(err)
			resizeImage(f.Name(), height, width)
		} else {
			fmt.Println(f.Name(), "isn't a picture, skipping.")
		}
	}
}

// SubImager acts as a way of gathering a sub-matrix of an image.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func resizeImage(imageName string, newHeight float64, newWidth float64) {
	r, _ := regexp.Compile(".jpe?g$")
	// open "test.jpg"
	file, err := os.Open(imageName)
	handleError(err)

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	handleError(err)
	file.Close()

	croppedImage := resize.Resize(uint(newWidth), 0, img, resize.Lanczos3)

	topCrop, err := smartcrop.Crop(croppedImage, int(newHeight), int(newWidth))
	handleError(err)
	fmt.Printf("Top crop: %+v\n", topCrop)

	sub, ok := croppedImage.(SubImager)
	if ok {

	} else {
		fmt.Println("No SubImage support")
		os.Exit(1)
	}

	resizedImage := sub.SubImage(topCrop)

	newImageName := r.ReplaceAllString(imageName, ".jpg")

	fmt.Println(imageName, "-->", newImageName)

	out, err := os.Create(strings.Join([]string{"../resized/", newImageName}, ""))
	handleError(err)
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, resizedImage, nil)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
