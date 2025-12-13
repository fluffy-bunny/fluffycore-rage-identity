package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Parse command line arguments
	dir := flag.String("dir", "./static-output", "Directory containing HTML files to fix")
	baseHref := flag.String("basehref", "management", "Base HREF value to replace with placeholder")
	title := flag.String("title", "MyApp", "Title value to replace with placeholder")
	flag.Parse()

	if *dir == "" {
		fmt.Println("Usage: fixup_static_html -dir=<directory>")
		os.Exit(1)
	}

	fmt.Printf("Fixing up HTML files in: %s\n", *dir)

	// Walk through directory and process all .html and app-worker.js files
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process HTML files
		if strings.HasSuffix(strings.ToLower(path), ".html") {
			fmt.Printf("Processing HTML: %s\n", path)
			return fixHTMLFile(path)
		}

		// Process app-worker.js files
		if filepath.Base(path) == "app-worker.js" {
			fmt.Printf("Processing Service Worker: %s\n", path)
			return fixServiceWorkerFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing files: %v\n", err)
		os.Exit(1)
	}

	// Create template files
	fmt.Println("Creating template files...")
	err = createTemplates(*dir, *baseHref, *title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating templates: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ HTML and Service Worker files fixed successfully!")
}

func fixHTMLFile(filePath string) error {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	// Regex to extract base href tag and move it to after <head>
	// Pattern: capture <head>, everything before base, base tag, everything after, </head>
	baseHrefRegex := regexp.MustCompile(`(?s)(<head[^>]*>)(.*?)(<base\s+href="[^"]+"\s*/?>)(.*?)(</head>)`)

	if baseHrefRegex.MatchString(contentStr) {
		// Move base href to right after <head> tag
		contentStr = baseHrefRegex.ReplaceAllString(contentStr, "$1\n    $3$2$4$5")
	}

	// Extract and move build version meta tags to just after <head> (or after base href if present)
	// Pattern: Find all meta tags with name="app-*"
	buildMetaRegex := regexp.MustCompile(`(?m)^\s*<meta\s+name="app-[^"]+"\s+content="[^"]+"\s*/?>.*?$\n?`)

	// Find all build meta tags
	buildMetas := buildMetaRegex.FindAllString(contentStr, -1)

	if len(buildMetas) > 0 {
		// Remove them from their current position
		contentStr = buildMetaRegex.ReplaceAllString(contentStr, "")

		// Insert them right after <head> tag (and base href if present)
		// First, find where to insert (after <head> or after <base> if it exists)
		insertRegex := regexp.MustCompile(`(<head[^>]*>\n(?:\s*<base[^>]+>\n)?)`)

		// Build the meta tags string
		metaTagsStr := strings.Join(buildMetas, "")

		// Insert the meta tags
		contentStr = insertRegex.ReplaceAllString(contentStr, "$1"+metaTagsStr)
	}

	// Remove entire <div class="wizard-container">...</div> and all its contents
	// This regex matches from the opening tag to the last closing </div> before <aside>
	wizardContainerRegex := regexp.MustCompile(`(?s)<div class="wizard-container">.*?</div>\s*(<aside)`)
	contentStr = wizardContainerRegex.ReplaceAllString(contentStr, "$1")

	// Write back to file
	err = os.WriteFile(filePath, []byte(contentStr), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func fixServiceWorkerFile(filePath string) error {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	// Find and fix the resourcesToCache array
	// Pattern: const resourcesToCache = [...];
	// We need to remove entries that end with "/" (root paths)
	resourcesToCacheRegex := regexp.MustCompile(`(const resourcesToCache = \[)([^\]]+)(\];)`)

	if resourcesToCacheRegex.MatchString(contentStr) {
		contentStr = resourcesToCacheRegex.ReplaceAllStringFunc(contentStr, func(match string) string {
			parts := resourcesToCacheRegex.FindStringSubmatch(match)
			if len(parts) != 4 {
				return match
			}

			prefix := parts[1]    // const resourcesToCache = [
			resources := parts[2] // the array contents
			suffix := parts[3]    // ];

			// Split by comma and filter out entries ending with "/"
			resourceArray := strings.Split(resources, ",")
			var filtered []string

			for _, resource := range resourceArray {
				resource = strings.TrimSpace(resource)
				if resource == "" {
					continue
				}

				// Check if this resource ends with "/" (inside quotes)
				// Pattern: "something/" or '/something/'
				endsWithSlashRegex := regexp.MustCompile(`^["'][^"']*\/["']$`)
				if !endsWithSlashRegex.MatchString(resource) {
					filtered = append(filtered, resource)
				}
			}

			// Rebuild the array
			return prefix + strings.Join(filtered, ",") + suffix
		})
	}

	// Also update the fetchWithCache function to never cache HTML
	// Check if it already has the HTML exclusion logic
	if !strings.Contains(contentStr, "request.destination === 'document'") {
		// Find the fetchWithCache function and add HTML exclusion
		fetchFuncRegex := regexp.MustCompile(`(?s)(async function fetchWithCache\(request\) \{\n)(.*?)(\n\})`)

		if fetchFuncRegex.MatchString(contentStr) {
			contentStr = fetchFuncRegex.ReplaceAllString(contentStr, `$1  // Never cache HTML files - always fetch fresh
  const url = new URL(request.url);
  if (request.destination === 'document' || 
      url.pathname.endsWith('.html') || 
      url.pathname.endsWith('/')) {
    return await fetch(request);
  }
  
$2$3`)
		}
	}

	// Write back to file
	err = os.WriteFile(filePath, []byte(contentStr), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
func createTemplates(dir, baseHref, title string) error {
	// Create index_template.html from index.html
	indexPath := filepath.Join(dir, "index.html")
	indexTemplatePath := filepath.Join(dir, "index_template.html")

	if _, err := os.Stat(indexPath); err == nil {
		content, err := os.ReadFile(indexPath)
		if err != nil {
			return fmt.Errorf("failed to read index.html: %w", err)
		}

		contentStr := string(content)

		// Replace baseHref with placeholder
		contentStr = strings.ReplaceAll(contentStr, "/"+baseHref+"/", "/{basehref}/")
		contentStr = strings.ReplaceAll(contentStr, "/"+baseHref, "/{basehref}")
		contentStr = strings.ReplaceAll(contentStr, baseHref, "{basehref}")

		// Replace title with placeholder
		titleRegex := regexp.MustCompile(`<title>[^<]*</title>`)
		contentStr = titleRegex.ReplaceAllString(contentStr, "<title>{title}</title>")

		err = os.WriteFile(indexTemplatePath, []byte(contentStr), 0644)
		if err != nil {
			return fmt.Errorf("failed to write index_template.html: %w", err)
		}
		fmt.Printf("✅ Created %s\n", indexTemplatePath)
	}

	// Create app_template.js from app.js
	appJsPath := filepath.Join(dir, "app.js")
	appTemplateJsPath := filepath.Join(dir, "app_template.js")

	if _, err := os.Stat(appJsPath); err == nil {
		content, err := os.ReadFile(appJsPath)
		if err != nil {
			return fmt.Errorf("failed to read app.js: %w", err)
		}

		contentStr := string(content)

		// Replace baseHref with placeholder (with and without slashes)
		contentStr = strings.ReplaceAll(contentStr, "/"+baseHref+"/", "/{basehref}/")
		contentStr = strings.ReplaceAll(contentStr, "/"+baseHref, "/{basehref}")
		contentStr = strings.ReplaceAll(contentStr, "\""+baseHref, "\"{basehref}")
		contentStr = strings.ReplaceAll(contentStr, baseHref, "{basehref}")

		err = os.WriteFile(appTemplateJsPath, []byte(contentStr), 0644)
		if err != nil {
			return fmt.Errorf("failed to write app_template.js: %w", err)
		}
		fmt.Printf("✅ Created %s\n", appTemplateJsPath)
	}

	return nil
}
