package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Parse command line arguments
	inputDir := flag.String("input", "./static_output", "Input directory containing static files")
	outputDir := flag.String("output", "./rebased_output", "Output directory for rebased static files")
	oldBase := flag.String("old-base", "/management", "Old base href path")
	newBase := flag.String("new-base", "/my/nested/app", "New base href path")
	flag.Parse()

	if *inputDir == "" || *outputDir == "" || *oldBase == "" || *newBase == "" {
		fmt.Println("Usage: rebase_static -input=<input_dir> -output=<output_dir> -old-base=<old_base> -new-base=<new_base>")
		fmt.Println("Example: rebase_static -input=./static_output -output=./rebased_output -old-base=/management -new-base=/my/nested/app")
		os.Exit(1)
	}

	// Ensure old base and new base have consistent format (with slashes)
	*oldBase = strings.TrimSuffix(*oldBase, "/")
	*newBase = strings.TrimSuffix(*newBase, "/")

	fmt.Printf("Rebasing static files from '%s' to '%s'\n", *oldBase, *newBase)
	fmt.Printf("Input directory: %s\n", *inputDir)
	fmt.Printf("Output directory: %s\n", *outputDir)

	// Remove output directory if it exists
	if err := os.RemoveAll(*outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing output directory: %v\n", err)
		os.Exit(1)
	}

	// Copy all files from input to output, processing text files
	err := filepath.Walk(*inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(*inputDir, path)
		if err != nil {
			return err
		}

		// Calculate output path
		outPath := filepath.Join(*outputDir, relPath)

		// If it's a directory, create it
		if info.IsDir() {
			return os.MkdirAll(outPath, 0755)
		}

		// Check if it's a text file that needs processing
		needsProcessing := false
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".html", ".js", ".webmanifest", ".json", ".css", ".txt":
			needsProcessing = true
		}

		if needsProcessing {
			return processTextFile(path, outPath, *oldBase, *newBase)
		}

		// For binary files, just copy
		return copyFile(path, outPath)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Static files rebased successfully!")
	fmt.Printf("Output directory: %s\n", *outputDir)
	fmt.Printf("\nTo test the rebased app, create a server pointing to: %s\n", *outputDir)
}

func processTextFile(inputPath, outputPath, oldBase, newBase string) error {
	// Read the file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", inputPath, err)
	}

	contentStr := string(content)

	// Replace old base with new base
	// Match patterns like /management/ and /management
	oldBaseWithSlash := oldBase + "/"
	newBaseWithSlash := newBase + "/"

	contentStr = strings.ReplaceAll(contentStr, oldBaseWithSlash, newBaseWithSlash)
	contentStr = strings.ReplaceAll(contentStr, oldBase, newBase)

	// Also handle quoted versions
	contentStr = strings.ReplaceAll(contentStr, `"`+oldBase+`"`, `"`+newBase+`"`)
	contentStr = strings.ReplaceAll(contentStr, `'`+oldBase+`'`, `'`+newBase+`'`)

	// Write to output file
	if err := os.WriteFile(outputPath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	fmt.Printf("Processed: %s\n", inputPath)
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
