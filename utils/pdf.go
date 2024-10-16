package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/jung-kurt/gofpdf"
)

type BarcodeInfo struct {
	Number string
	Label  string
}

func GenerateBarcodePDF(barcodes []BarcodeInfo, outputPath string) error {
	// PDF settings
	pageWidth, _ := 210.0, 297.0
	margin := 10.0
	barcodeWidth, barcodeHeight := (pageWidth-margin*2-10.0)/2.0, 30.0
	labelHeight := 10.0
	horizontalSpacing, verticalSpacing := pageWidth-margin*2-barcodeWidth*2, 5.0

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("Roboto", "", "assets/Roboto-Regular.ttf")
	pdf.AddUTF8Font("Roboto-Bold", "", "assets/Roboto-Bold.ttf")
	pdf.SetMargins(margin, margin, margin)

	itemsPerPage := 10 // 2 columns * 5 rows
	for i, info := range barcodes {
		if i%itemsPerPage == 0 {
			pdf.AddPage()

			// Add barcode to PDF
			pdf.Image("assets/jamkat.png", margin+5, margin, barcodeWidth, barcodeWidth*(241.0/577.0), false, "", 0, "")

		}

		// Calculate position
		col := (i % 2) * int(barcodeWidth+horizontalSpacing)
		row := ((i % itemsPerPage) / 2) * int(barcodeHeight+labelHeight+verticalSpacing)
		x := float64(col) + margin
		y := float64(row) + margin + barcodeWidth*(241.0/577.0) + 5

		// Generate barcode
		barcodeImg, err := generateBarcode(info.Number, int(barcodeWidth*5), int(barcodeHeight*5))
		if err != nil {
			return fmt.Errorf("error generating barcode: %v", err)
		}

		// Convert to 8-bit and get PNG bytes
		pngBytes, err := convertToEightBitPNG(barcodeImg)
		if err != nil {
			return fmt.Errorf("error converting barcode to PNG: %v", err)
		}

		// Add barcode to PDF
		imgPath := fmt.Sprintf("temp_barcode_%d.png", i)
		if err := os.WriteFile(imgPath, pngBytes, 0644); err != nil {
			return fmt.Errorf("error saving barcode image: %v", err)
		}
		pdf.Image(imgPath, x, y, barcodeWidth, barcodeHeight, false, "", 0, "")
		os.Remove(imgPath) // Clean up temporary file

		// Add label
		pdf.SetFont("Roboto-Bold", "", 24) // Increased font size to 12 and made it bold
		labelWidth := pdf.GetStringWidth(info.Label)
		labelX := x + (barcodeWidth-labelWidth)/2 // Center the label
		pdf.Text(labelX, y+barcodeHeight+10, info.Label)

	}

	return pdf.OutputFileAndClose(outputPath)
}

func generateBarcode(number string, width, height int) (barcode.Barcode, error) {
	barcodeImg, err := code128.Encode(number)
	if err != nil {
		return nil, err
	}

	// Scale the barcode
	scaledImg, err := barcode.Scale(barcodeImg, width, height)
	if err != nil {
		return nil, err
	}

	return scaledImg, nil
}
func convertToEightBitPNG(img barcode.Barcode) ([]byte, error) {
	// Create a new 8-bit image
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Copy the image data
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, newImg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
