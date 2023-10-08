package src

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// determine the file extension based on the content type
func extractImageExtension(contentType string) (string, error) {
	if contentType == "" {
		return "", errors.New("no content type provided")
	}

	splittedContentType := strings.Split(contentType, ";")
	mainType := splittedContentType[0]

	switch mainType {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	default:
		return "", fmt.Errorf("unsupported image type: %s", contentType)
	}
}

// save the image content to the specified path
func saveImage(imageContent io.Reader, filePathWithName string) error {
	file, err := os.Create(filePathWithName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, imageContent); err != nil {
		return err
	}

	return nil
}

// findDeepestDir recursively finds the deepest directory starting from baseDir.
func findDeepestDir(baseDir string) (string, error) {
	var deepestDir string
	var deepestDirLevel int

	err := filepath.WalkDir(baseDir, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			currentLevel := len(strings.Split(currentPath, string(os.PathSeparator)))
			if currentLevel > deepestDirLevel {
				deepestDir = currentPath
				deepestDirLevel = currentLevel
			}
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return deepestDir, nil
}

// GetPoster downloads an image from the provided URL and saves it to the specified file path.
func GetPoster(url, filePath string) error {
	// Make a GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response is OK
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to get the image")
	}

	// Ensure the content type is an image
	contentType := resp.Header.Get("Content-Type")
	ext, err := extractImageExtension(contentType)
	if err != nil {
		return err
	}

	// Find the deepest directory to save the poster
	deepestDir, err := findDeepestDir(filePath)
	if err != nil {
		return fmt.Errorf("error finding deepest directory: %w", err)
	}

	filePathWithName := path.Join(deepestDir, "poster"+ext)
	if err := saveImage(resp.Body, filePathWithName); err != nil {
		return fmt.Errorf("error saving image: %w", err)
	}

	return nil
}
