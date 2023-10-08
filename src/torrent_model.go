package src

import (
	"errors"
)

// TorrentRequest represents the expected request payload.
type TorrentRequest struct {
	Magnet    string `json:"magnet"`
	Folder    string `json:"folder"`
	Name      string `json:"name"`
	PosterURL string `json:"posterUrl"`
}

func formatErrors(errs []error) error {
	errStr := "Invalid request: \n"
	for _, err := range errs {
		errStr += err.Error() + "\n"
	}
	return errors.New(errStr)
}

// Validate checks the fields of TorrentRequest for any invalid or missing data.
func (tr *TorrentRequest) Validate() error {
	errs := []error{}
	if tr.Magnet == "" {
		errs = append(errs, errors.New("magnet is required"))
	}
	if tr.Folder == "" {
		errs = append(errs, errors.New("folder is required"))
	}
	if tr.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}
	if tr.PosterURL == "" {
		errs = append(errs, errors.New("posterUrl is required"))
	}
	if len(errs) > 0 {
		return formatErrors(errs)
	}
	return nil
}
