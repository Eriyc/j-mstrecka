package pdf

import (
	"fmt"
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
	//	pageWidth, pageHeight := 210.0, 297.0
	margin := 10.0
	barcodeWidth, barcodeHeight := 80.0, 40.0
	labelHeight := 10.0
	spacing := 5.0

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(margin, margin, margin)

	itemsPerPage := 10 // 2 columns * 5 rows
	for i, info := range barcodes {
		if i%itemsPerPage == 0 {
			pdf.AddPage()
		}

		// Calculate position
		col := (i % 2) * int(barcodeWidth+spacing)
		row := ((i % itemsPerPage) / 2) * int(barcodeHeight+labelHeight+spacing)
		x := float64(col) + margin
		y := float64(row) + margin

		// Generate barcode
		barcodeImg, err := generateBarcode(info.Number)
		if err != nil {
			return fmt.Errorf("error generating barcode: %v", err)
		}

		// Add barcode to PDF
		imgPath := fmt.Sprintf("temp_barcode_%d.png", i)
		if err := saveBarcodePNG(barcodeImg, imgPath); err != nil {
			return fmt.Errorf("error saving barcode image: %v", err)
		}
		pdf.Image(imgPath, x, y, barcodeWidth, barcodeHeight, false, "", 0, "")
		os.Remove(imgPath) // Clean up temporary file

		// Add label
		pdf.SetFont("Arial", "", 10)
		pdf.Text(x, y+barcodeHeight+5, info.Label)
	}

	return pdf.OutputFileAndClose(outputPath)
}

func generateBarcode(number string) (barcode.Barcode, error) {
	return code128.Encode(number)
}

func saveBarcodePNG(barcodeImg barcode.Barcode, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, barcodeImg)
}
