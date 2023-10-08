package src

import (
	"errors"
)

type SetPosterRequest struct {
	Name   string `json:"name"`
	Folder string `json:"folder"`
	URL    string `json:"url"`
}

func formatPosterErrors(errs []error) error {
	errStr := "Invalid request: \n"
	for _, err := range errs {
		errStr += err.Error() + "\n"
	}
	return errors.New(errStr)
}

// Validate checks the fields of SetPosterRequest for any invalid or missing data.
func (tr *SetPosterRequest) Validate() error {
	errs := []error{}
	if tr.Folder == "" {
		errs = append(errs, errors.New("folder is required"))
	}
	if tr.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}
	if tr.URL == "" {
		errs = append(errs, errors.New("url is required"))
	}
	if len(errs) > 0 {
		return formatErrors(errs)
	}
	return nil
}
