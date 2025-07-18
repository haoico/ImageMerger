package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	mergeEqualSizedImages()
	mergeDifferentSizedImages()
}

// Function to merge equal sized images
func mergeEqualSizedImages() {
	// Define the image file paths
	images := []string{"pt1.jpg", "pt2.jpg", "pt3.jpg", "pt4.jpg", "pt5.jpg", "pt6.jpg", "pt7.jpg", "pt8.jpg", "pt9.jpg", "pt10.jpg", "pt11.jpg", "pt12.jpg"}

	// Open the images
	var imgList []image.Image
	var totalHeight, maxWidth int

	for _, file := range images {
		img, err := openImage(file)
		if err != nil {
			log.Fatal(err)
		}
		imgList = append(imgList, img)
		totalHeight += img.Bounds().Dy() // Add the width of the current image to the total width
		if img.Bounds().Dx() > maxWidth {
			maxWidth = img.Bounds().Dx() // Find the maximum height
		}
	}

	// Create a new blank image with combined width and maximum height
	finalImg := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))

	// Paste the images into the final image horizontally
	yOffset := 0
	for _, img := range imgList {
		draw.Draw(finalImg, img.Bounds().Add(image.Pt(0, yOffset)), img, image.Point{}, draw.Over)
		yOffset += img.Bounds().Dy()
	}

	// Save the combined image as a JPEG
	combinedImageFile := "merged_image.jpg"
	saveImageAsJPEG(combinedImageFile, finalImg)

	// Convert the combined image to a PDF
	pdfFile := "output.pdf"
	convertImageToPDF(combinedImageFile, pdfFile)

	fmt.Println("Process completed successfully!")
}

// Function to merge equal sized images
func mergeDifferentSizedImages() {
	// Define the image file paths
	equalImages := []string{"pt1.jpg", "pt2.jpg", "pt3.jpg", "pt4.jpg", "pt5.jpg", "pt6.jpg", "pt7.jpg", "pt8.jpg", "pt9.jpg"}
	unequalImages := []string{"pt10.jpg", "pt11.jpg", "pt12.jpg", "pt13.jpg", "pt14.jpg", "pt15.jpg"}

	// Open the images
	var equalImgList []image.Image
	var unequalImgList []image.Image
	var totalHeight, maxWidth, unequalWidth int

	for _, file := range equalImages {
		img, err := openImage(file)
		if err != nil {
			log.Fatal(err)
		}
		equalImgList = append(equalImgList, img)
		totalHeight += img.Bounds().Dy() // Add the width of the current image to the total width
		if img.Bounds().Dx() > maxWidth {
			maxWidth = img.Bounds().Dx() // Find the maximum height
		}
	}

	for _, file := range unequalImages {
		img, err := openImage(file)
		if err != nil {
			log.Fatal(err)
		}
		unequalImgList = append(unequalImgList, img)
		totalHeight += (img.Bounds().Dy() / 2) // Add the width of the current image to the total width
		if img.Bounds().Dx() > unequalWidth {
			unequalWidth = img.Bounds().Dx() // Find the maximum height
		}
	}

	// Create a new blank image with combined width and maximum height
	finalImg := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))

	// Paste the images into the final image horizontally
	yOffset := 0
	xOffset := 66

	for _, img := range equalImgList {
		draw.Draw(finalImg, img.Bounds().Add(image.Pt(0, yOffset)), img, image.Point{}, draw.Over)
		yOffset += img.Bounds().Dy()
	}

	imgCounter := 0

	for _, img := range unequalImgList {
		draw.Draw(finalImg, img.Bounds().Add(image.Pt(xOffset, yOffset)), img, image.Point{}, draw.Over)
		if imgCounter == 1 {
			yOffset += img.Bounds().Dy()
			xOffset = 66
			imgCounter = 0
		} else {
			xOffset = unequalWidth
			imgCounter = 1
		}
	}

	// Save the combined image as a JPEG
	combinedImageFile := "merged_image.jpg"
	saveImageAsJPEG(combinedImageFile, finalImg)

	// Convert the combined image to a PDF
	pdfFile := "output.pdf"
	convertImageToPDF(combinedImageFile, pdfFile)

	fmt.Println("Process completed successfully!")
}

// Function to open an image
func openImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open image %s: %w", filePath, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("unable to decode image %s: %w", filePath, err)
	}
	return img, nil
}

// Function to save an image as JPEG
func saveImageAsJPEG(filePath string, img image.Image) {
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("failed to create file %s: %v", filePath, err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, img, nil)
	if err != nil {
		log.Fatalf("failed to save image %s: %v", filePath, err)
	}
}

// Function to convert an image to a PDF
func convertImageToPDF(imagePath, pdfPath string) {
	// Initialize PDF
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("failed to open image %s: %v", imagePath, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image %s: %v", imagePath, err)
	}

	// Set PDF properties
	pdf.AddPage()

	// Convert the image size to fit the page
	pageWidth, pageHeight := pdf.GetPageSize()
	imgWidth := float64(img.Bounds().Dx())
	imgHeight := float64(img.Bounds().Dy())

	// Scale image to fit within the A4 page
	scale := min(pageWidth/imgWidth, pageHeight/imgHeight)
	imgWidth *= scale
	imgHeight *= scale

	// Add image to the PDF
	pdf.Image(imagePath, 0, 0, imgWidth, imgHeight, false, "", 0, "")

	// Output the PDF to a file
	err = pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		log.Fatalf("failed to save PDF %s: %v", pdfPath, err)
	}
}

// Helper function to calculate the minimum of two float64 numbers
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
