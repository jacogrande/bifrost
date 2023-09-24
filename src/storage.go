package src

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

func saveMapToFile(folderMap map[string]string, filePath string) error {
	// Convert map to JSON
	data, err := json.Marshal(folderMap)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0644)
	return err
}

func loadMapFromFile(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// parse json
	var folderMap map[string]string
	if len(data) == 0 {
		folderMap = make(map[string]string)
		return folderMap, nil
	}
	err = json.Unmarshal(data, &folderMap)
	if err != nil {
		return nil, err
	}

	return folderMap, nil
}

func getStoragePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home path: %s", err)
	}
	bifrostDir := filepath.Join(homeDir, ".bifrost")
	return filepath.Join(bifrostDir, "folders.json")
}

// creates the .bifrost directory if it doesn't exist
func ensureHomeDir() {
	// ensure that the .bifrost directory exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %s", err)
	}
	// make the .bifrost directory if it doesn't exist
	bifrostDir := filepath.Join(homeDir, ".bifrost")
	if _, err := os.Stat(bifrostDir); os.IsNotExist(err) {
		os.Mkdir(bifrostDir, 0755)
	}
}

// creates the .bifrost/folders.json file if it doesn't exist
func ensureFoldersFile() {
	path := getStoragePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create the file if it does not exist
		_, err := os.Create(path)
		if err != nil {
			log.Fatalf("Failed creating file: %s", err)
		}
	}
}

func SaveFolderInit(name string, path string) error {
	ensureHomeDir()
	ensureFoldersFile()
	storagePath := getStoragePath()
	log.Println(storagePath)
	folderMap, loadErr := loadMapFromFile(storagePath)
	if loadErr != nil {
		return loadErr
	}
	_, folderExists := folderMap[name]
	if folderExists {
		return errors.New("folder already exists")
	}
	folderMap[name] = path
	saveErr := saveMapToFile(folderMap, storagePath)
	if saveErr != nil {
		return saveErr
	}

	log.Printf("Saved folder '%s' with path '%s'\n", name, path)

	return nil
}

func GetFolderPath(folderName string, torrentName string) (string, error) {
	storagePath := getStoragePath()
	folderMap, loadErr := loadMapFromFile(storagePath)
	if loadErr != nil {
		return "", loadErr
	}
	path, folderExists := folderMap[folderName]
	if !folderExists {
		return "", errors.New("folder does not exist")
	}
	path = filepath.Join(path, torrentName)
	return path, nil
}
