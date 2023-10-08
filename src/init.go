package src

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func promptUser(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input) // Remove potential whitespace or newline characters.
	return input, nil
}

func getFolderPath(reader *bufio.Reader, folderName string) (string, error) {
	return promptUser(reader, fmt.Sprintf("Enter your path for the folder '%s': ", folderName))
}

func getFolderName(reader *bufio.Reader) (string, error) {
	return promptUser(reader, "Enter your folder name: ")
}

func checkPermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Extract permission bits
	perm := info.Mode().Perm()

	// Check for read and write permissions for the user
	if perm&(1<<7) == 0 || perm&(1<<8) == 0 {
		return errors.New("folder does not have read and write permissions")
	}

	return nil
}

func Init() {
	reader := bufio.NewReader(os.Stdin)
	// infinite loop
	for {
		// get folder data
		folderName, err := getFolderName(reader)
		if err != nil {
			fmt.Println("Error getting folder name:", err)
			return
		}

		folderPath, err := getFolderPath(reader, folderName)
		if err != nil {
			fmt.Println("Error getting folder path:", err)
			return
		}
		// check permissions
		permissionsErr := checkPermissions(folderPath)
		if permissionsErr != nil {
			fmt.Println("Error checking permissions:", permissionsErr)
		} else {
			// save folder
			saveErr := SaveFolderInit(folderName, folderPath)
			if saveErr != nil {
				fmt.Println("Error saving folder:", saveErr)
			}
		}

		fmt.Print("Would you like to add another folder? (y/n): ")
		choice, _ := reader.ReadString('\n') // For simplicity, we're ignoring the error here.
		choice = strings.TrimSpace(strings.ToLower(choice))

		if choice != "y" {
			break
		}
	}
}
