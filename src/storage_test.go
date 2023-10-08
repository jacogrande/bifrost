package src

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveMapToFile(t *testing.T) {
	// Setup
	tmpfile, err := os.CreateTemp("", "example.*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	folderMap := map[string]string{
		"test": "path/to/test",
	}

	// Test
	err = saveMapToFile(folderMap, tmpfile.Name())
	assert.NoError(t, err)

	// Check file content
	content, err := os.ReadFile(tmpfile.Name())
	assert.NoError(t, err)

	var readFolderMap map[string]string
	err = json.Unmarshal(content, &readFolderMap)
	assert.NoError(t, err)

	assert.Equal(t, folderMap, readFolderMap)
}

func TestLoadMapFromFile(t *testing.T) {
	// Setup
	tmpfile, err := os.CreateTemp("", "example.*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	// initialize folder map
	folderMap := map[string]string{
		"test": "path/to/test",
	}
	data, err := json.Marshal(folderMap)
	assert.NoError(t, err)

	// write folder file
	err = os.WriteFile(tmpfile.Name(), data, 0644)
	assert.NoError(t, err)

	// Test
	readFolderMap, err := loadMapFromFile(tmpfile.Name())
	assert.NoError(t, err)

	assert.Equal(t, folderMap, readFolderMap)
}
