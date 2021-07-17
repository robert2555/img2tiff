package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/image/tiff"
)

func loadImg(imgPath string) (img image.Image, err error) {
	// Open Image File
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return img, err // Catch err and exit
	}
	// Get Image Type
	_, imgType, _ := image.Decode(imgFile)
	fmt.Printf("Found file with Type: %v \n", imgType)
	// Resets the begining of the file for the next read
	imgFile.Seek(0, 0)
	// Read image into image var
	switch imgType {
	case "png":
		img, err = png.Decode(imgFile)
		if err != nil {
			return img, err // Catch err and exit
		}
	case "jpeg":
		img, err = jpeg.Decode(imgFile)
		if err != nil {
			return img, err // Catch err and exit
		}
	case "gif":
		img, err = gif.Decode(imgFile)
		if err != nil {
			return img, err // Catch err and exit
		}
	default:
		fmt.Printf("%v is not a valid image type!", imgType)
		os.Exit(1)
	}

	return img, nil
}

func getContentType(file *os.File) string {
	// Create a buffer with 512 Byte
	buffer := make([]byte, 512)
	// Read 512Byte(buffer) from the file and write it in the byte slice named buffer
	len, err := file.Read(buffer)
	if err != nil {
		fmt.Println("Error Read from file", err)
	}
	// Get the content type from the 512 bytes of the file
	ct := mimetype.Detect(buffer[:len])

	return ct.String()
}

func main() {
	// Get the commandline parameters
	var imgPath, tiffPath string

	switch len(os.Args) {
	case 2:
		imgPath = os.Args[1]
		idx := strings.Index(imgPath, ".")
		tiffPath = imgPath[:idx] + ".tiff"
	case 3:
		imgPath = os.Args[1]
		tiffPath = os.Args[2]
	default:
		fmt.Printf("Usage: %v <Image File> <New File>\n", os.Args[0])
		os.Exit(1)
	}

	// Load the image
	img, err := loadImg(imgPath)
	if err != nil {
		fmt.Println(err)
	}

	// Create a new file
	tiffFile, err := os.Create(tiffPath)
	if err != nil {
		fmt.Println(err)
	}

	// Convert img to tiff and write into new created file
	fmt.Println("Decoding image...")
	err = tiff.Encode(tiffFile, img, nil)
	if err != nil {
		fmt.Println(err)
	}

	// Close the new file and open it for file inspection
	tiffFile.Close()
	tiffFile, err = os.Open(tiffPath)

	// Get the filetype, then close the file
	fileType := getContentType(tiffFile)
	tiffFile.Close()

	// Check if the tiff is valid
	if fileType == "image/tiff" {
		fmt.Println("Success! New file is a valid tiff!")
	} else {
		fmt.Println("Tiff file content is not tiff! Conversion maybe failed.")
		os.Exit(1)
	}
}
