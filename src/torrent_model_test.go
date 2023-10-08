package src

import (
	"testing"
)

func TestTorrentRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request TorrentRequest
		wantErr bool
	}{
		{
			name:    "Valid Request",
			request: TorrentRequest{Magnet: "magnet:?xt=urn:btih:example", Folder: "folder", Name: "name", PosterURL: "http://example.com/poster.jpg"},
			wantErr: false,
		},
		{
			name:    "Empty Magnet",
			request: TorrentRequest{Magnet: "", Folder: "folder", Name: "name", PosterURL: "http://example.com/poster.jpg"},
			wantErr: true,
		},
		{
			name:    "Empty Folder",
			request: TorrentRequest{Magnet: "magnet:?xt=urn:btih:example", Folder: "", Name: "name", PosterURL: "http://example.com/poster.jpg"},
			wantErr: true,
		},
		{
			name:    "Empty Name",
			request: TorrentRequest{Magnet: "magnet:?xt=urn:btih:example", Folder: "folder", Name: "", PosterURL: "http://example.com/poster.jpg"},
			wantErr: true,
		},
		{
			name:    "Empty PosterURL",
			request: TorrentRequest{Magnet: "magnet:?xt=urn:btih:example", Folder: "folder", Name: "name", PosterURL: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("TorrentRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
