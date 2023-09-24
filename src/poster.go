package src

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

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
	if contentType != "image/jpeg" && contentType != "image/png" {
		return errors.New("URL does not point to a valid image")
	}

	// Determine the extension for the image
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		return errors.New("unsupported image type")
	}

	// Open a file for writing
	file, err := os.Create(path.Join(filePath, name+ext))
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the image content to the file
	_, err = io.Copy(file, resp.Body)
	return err
}
