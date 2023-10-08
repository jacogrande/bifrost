package src

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
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

// GetPoster downloads an image from the provided URL and saves it to the specified file path.
func GetPoster(url, filePath, name string) error {
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

	filePathWithName := path.Join(filePath, name+ext)
	if err := saveImage(resp.Body, filePathWithName); err != nil {
		return fmt.Errorf("error saving image: %w", err)
	}

	return nil
}
