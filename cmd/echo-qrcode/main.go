package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
	qrcode "github.com/skip2/go-qrcode"
	qrsvg "github.com/wamuir/svg-qr-code"
	gotp "github.com/xlzd/gotp"
)

const (
	secret    = "4S62BZNFXXSZLCRO"
	indexHTML = `
	<!DOCTYPE html>
	<html>
	<head>
		<title>QR Code</title>
		<script>
			function verify() {
				var input = document.getElementById("verifyField").value;
				fetch('/verify?input=' + input)
					.then(response => response.text())
					.then(data => alert(data));
			}
		</script>
	</head>
	<body>
		<div style="width: 200px; height: 200px; display: flex; justify-content: center; align-items: center;">
			<img src="data:image/png;base64,{{PNG}}" alt="QR Code" style="max-width: 100%; max-height: 100%;" /> 
		</div>
		<br/>
		<div style="display: flex; flex-direction: column; gap: 10px;">
			<input type="text" id="verifyField" placeholder="Enter text to verify" />
			<button onclick="verify()">Verify</button>
		</div>
	</body>
	</html>
`
)

func generateQRCode(c echo.Context) error {
	// Get the text you want to encode (e.g., a URL)
	//text := "https://example.org"
	otp := gotp.NewDefaultTOTP(secret)

	provisioningUri := otp.ProvisioningUri("a@b.com", "rage.identity")

	_, err := qrsvg.New(provisioningUri)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error generating QR code")
	}
	//fmt.Println(qr.String())
	var pngB []byte
	pngB, _ = qrcode.Encode(provisioningUri, qrcode.Medium, 256)
	base64Str := base64.StdEncoding.EncodeToString(pngB)

	// set svg content type
	//c.Response().Header().Set("Content-Type", "image/svg+xml")
	//svg := qr.SVG().String()
	dd := strings.Replace(indexHTML, "{{PNG}}", base64Str, -1)
	return c.HTML(http.StatusOK, dd)
	// Write the image to the response
	/*
		_, err = c.Response().Write([]byte(svg))
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error writing QR code image")
		}
	*/

	//return nil
}

func main() {
	e := echo.New()

	// Route to generate QR code
	e.GET("/qrcode", generateQRCode)
	e.GET("/verify", func(c echo.Context) error {
		input := c.QueryParam("input")
		// Add your verification logic here
		otp := gotp.NewDefaultTOTP(secret)
		verified := otp.Verify(input, time.Now().Unix())
		return c.String(http.StatusOK, fmt.Sprintf("Verified: %s, %v", input, verified))
	})
	// Start the server
	e.Start(":9055")
}
