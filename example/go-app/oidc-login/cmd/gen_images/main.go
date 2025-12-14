package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func main() {
	// Determine paths relative to the gen_images directory
	// gen_images is at: example/go-app/oidc-login/cmd/gen_images
	// SVG source is at: example/go-app/oidc-login/web/m_logo.svg
	// Output should go to: example/go-app/oidc-login/web/

	svgPath := filepath.Join("..", "..", "web", "m_logo.svg")
	outputDir := filepath.Join("..", "..", "web")

	// Read the SVG file
	svgData, err := os.ReadFile(svgPath)
	if err != nil {
		log.Fatalf("Error reading SVG file %s: %v", svgPath, err)
	}

	// Define the sizes for different favicon formats
	sizes := map[string]int{
		"favicon-16x16.png":          16,
		"favicon-32x32.png":          32,
		"favicon-64x64.png":          64,
		"apple-touch-icon.png":       180,
		"android-chrome-192x192.png": 192,
		"android-chrome-512x512.png": 512,
		"mstile-150x150.png":         150,
	}

	fmt.Printf("Converting SVG to PNG and generating favicons...\n")
	fmt.Printf("Source: %s\n", svgPath)
	fmt.Printf("Output: %s\n\n", outputDir)

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	for filename, size := range sizes {
		outputPath := filepath.Join(outputDir, filename)
		err := convertSVGToPNG(svgData, size, outputPath)
		if err != nil {
			log.Printf("Error converting to %s (%dx%d): %v", filename, size, size, err)
			continue
		}
		fmt.Printf("✓ Generated %s (%dx%d)\n", filename, size, size)
	}

	fmt.Println("\n✅ Conversion complete!")
}

func convertSVGToPNG(svgData []byte, size int, outputFile string) error {
	// Parse the SVG
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData))
	if err != nil {
		return fmt.Errorf("error parsing SVG: %v", err)
	}

	// Set the target size
	icon.SetTarget(0, 0, float64(size), float64(size))

	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Create a scanner and rasterize the SVG
	scanner := rasterx.NewScannerGV(size, size, img, img.Bounds())
	raster := rasterx.NewDasher(size, size, scanner)

	// Draw the SVG
	icon.Draw(raster, 1.0)

	// Save the PNG file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		return fmt.Errorf("error encoding PNG: %v", err)
	}

	return nil
}
