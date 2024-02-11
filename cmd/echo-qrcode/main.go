package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
	qrsvg "github.com/wamuir/svg-qr-code"
)

func generateQRCode(c echo.Context) error {
	// Get the text you want to encode (e.g., a URL)
	text := "https://example.org"

	qr, err := qrsvg.New(text)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error generating QR code")
	}
	//fmt.Println(qr.String())

	// set svg content type
	c.Response().Header().Set("Content-Type", "image/svg+xml")
	svg := qr.String()
	// Write the image to the response
	_, err = c.Response().Write([]byte(svg))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error writing QR code image")
	}

	// Generate a QR code image (256x256 pixels, medium error recovery level)
	_, err = qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error generating QR code")
	}

	// Set the response header for the image
	//	c.Response().Header().Set("Content-Type", "image/png")

	// Write the image to the response
	//	_, err = c.Response().Write(png)
	//	if err != nil {
	//		return c.String(http.StatusInternalServerError, "Error writing QR code image")
	//	}

	return nil
}

func main() {
	e := echo.New()

	// Route to generate QR code
	e.GET("/qrcode", generateQRCode)

	// Start the server
	e.Start(":9055")
}
