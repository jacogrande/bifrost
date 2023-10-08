package src

import "testing"

func runTestCase(t *testing.T, input string, expected string, hasErr bool) {
	output, err := extractImageExtension(input)
	if hasErr {
		if err == nil {
			t.Errorf("Expected an error for input: %s", input)
		}
	} else {
		if err != nil {
			t.Errorf("Expected no error for input: %s", input)
		}
		if output != expected {
			t.Errorf("Expected %s for input: %s", expected, input)
		}
	}
}

func TestExtractImageExtension(t *testing.T) {
	// define a slice of tests
	tests := []struct {
		input    string
		expected string
		hasErr   bool
	}{
		{"image/jpeg", ".jpg", false},
		{"image/png", ".png", false},
		{"image/jpeg; charset=UTF-8", ".jpg", false},
		{"image/gif", "", true},
	}

	for _, test := range tests {
		runTestCase(t, test.input, test.expected, test.hasErr)
	}
}
