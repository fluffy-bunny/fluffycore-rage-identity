package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Parse command line arguments
	staticDir := flag.String("dir", "./static_output", "Directory containing static files")
	basePath := flag.String("base", "/management", "Base path for serving the app")
	port := flag.String("port", "3001", "Port to listen on")
	flag.Parse()

	// Ensure basePath starts with / and doesn't end with /
	*basePath = "/" + strings.Trim(*basePath, "/")

	fmt.Printf("Starting server...\n")
	fmt.Printf("Static directory: %s\n", *staticDir)
	fmt.Printf("Base path: %s\n", *basePath)
	fmt.Printf("Port: %s\n", *port)
	fmt.Printf("\nAccess the app at: http://localhost:%s%s/\n", *port, *basePath)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Root page with selector to the base path
	e.GET("/", func(c echo.Context) error {
		html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Authentication Portal</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #0a2540 0%%, #1a365d 100%%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 60px 40px;
            max-width: 500px;
            width: 100%%;
            text-align: center;
        }
        h1 {
            color: #0a2540;
            font-size: 32px;
            font-weight: 700;
            margin-bottom: 16px;
        }
        p {
            color: #666;
            font-size: 16px;
            line-height: 1.6;
            margin-bottom: 40px;
        }
        .btn {
            display: inline-block;
            background: #00C896;
            color: white;
            text-decoration: none;
            padding: 16px 48px;
            border-radius: 8px;
            font-size: 18px;
            font-weight: 600;
            transition: all 0.2s ease;
            box-shadow: 0 4px 12px rgba(0, 200, 150, 0.3);
        }
        .btn:hover {
            background: #00b587;
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(0, 200, 150, 0.4);
        }
        .btn:active {
            transform: translateY(0);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Authentication Portal</h1>
        <p>Click below to access the management portal</p>
        <a href="%s/" class="btn">Continue to Management</a>
    </div>
</body>
</html>`, *basePath)
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "0")
		return c.HTML(http.StatusOK, html)
	})

	// Serve static files from the specified directory at the base path
	e.Static(*basePath, *staticDir)

	// Middleware to set cache headers and Content-Type AFTER static file handler
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			path := c.Request().URL.Path

			// Set Content-Type for .wasm files
			if strings.HasSuffix(path, ".wasm") {
				c.Response().Header().Set("Content-Type", "application/wasm")
			}

			// No cache for HTML files - always fetch fresh
			if strings.HasSuffix(path, ".html") || path == *basePath+"/" || path == *basePath || path == "/" {
				c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				c.Response().Header().Set("Pragma", "no-cache")
				c.Response().Header().Set("Expires", "0")
			} else if strings.HasPrefix(path, *basePath+"/") {
				// Cache assets with version query params for 1 year
				c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}

			return err
		}
	})

	// SPA fallback for non-file routes (must come after Static)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				httpErr, ok := err.(*echo.HTTPError)
				if ok && httpErr.Code == 404 && strings.HasPrefix(c.Request().URL.Path, *basePath+"/") {
					// Return index.html for 404s under the base path to support SPA routing
					return c.File(*staticDir + "/index.html")
				}
			}
			return err
		}
	})

	e.Logger.Fatal(e.Start(":" + *port))
}
