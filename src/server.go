package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

const TIMEOUT_DURATION = 30 * time.Minute

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/download", downloadHandler)
	mux.HandleFunc("/setPoster", setPosterHandler)
	log.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", mux)
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	log.Println("Error: ", err)
	http.Error(w, err.Error(), statusCode)
}

// expected request body
// magnet: string
// folder: string
// name: string
// posterUrl: string
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// verify request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// validate the request body
	var req TorrentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	folderPath, err := GetFolderPath(req.Folder, req.Name)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	// async channels
	errCh := make(chan error)
	successCh := make(chan bool)

	// start download
	go func() {
		if err := downloadTorrent(req.Magnet, folderPath, errCh, successCh); err != nil {
			errCh <- err
		}
	}()

	// wait for download to complete or timeout
	select {
	case err := <-errCh:
		writeError(w, err, http.StatusInternalServerError)
		return
	case <-successCh:
		if err := GetPoster(req.PosterURL, folderPath); err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Download completed successfully!"))
			return
		}
	case <-time.After(TIMEOUT_DURATION):
		w.Write([]byte("Download initiated or taking longer than expected"))
		return
	}
}

func handleWriteError(err error, errCh chan<- error) {
	log.Println(err)
	errCh <- err
}

func getDownloadStatus(t *torrent.Torrent) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	progressBarWidth := 50

	// print the torrent stats
	stats := t.Stats()
	fmt.Printf("Downloading %s\n", t.Name())
	fmt.Printf("Peers: %d\n", stats.ActivePeers)
	fmt.Printf("Seeders: %d\n", stats.ConnectedSeeders)

	for {
		select {
		case <-ticker.C:
			bytesCompleted := t.BytesCompleted()
			totalLength := t.Length()
			percentComplete := float64(bytesCompleted) / float64(totalLength) * 100
			progress := int(float64(progressBarWidth) * percentComplete / 100)
			progressBar := fmt.Sprintf("[%s%s]", strings.Repeat("=", progress), strings.Repeat(" ", progressBarWidth-progress))

			// Use ANSI escape codes to go back to the beginning of the line and clear it.
			// "\r" returns to the beginning of the line.
			// "\033[K" clears from the cursor to the end of the line.
			fmt.Printf("\r\033[KDownload status: %.2f%% complete %s", percentComplete, progressBar)
		case <-t.Closed():
			fmt.Printf("\r\033[KDownload status: 100%% complete [%s]", strings.Repeat("=", progressBarWidth))
			// Exit the goroutine when the torrent is closed
			return
		}
	}
}

func downloadTorrent(magnet string, saveDirectory string, errCh chan<- error, successCh chan<- bool) error {
	// initialize torrent client
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = saveDirectory
	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	// add torrent
	t, err := client.AddMagnet(magnet)
	if err != nil {
		return err
	}

	// wait for torrent to be ready
	<-t.GotInfo()
	t.SetOnWriteChunkError(func(err error) {
		handleWriteError(err, errCh)
	})

	go getDownloadStatus(t)

	t.DownloadAll()

	// wait for download to complete
	if !client.WaitAll() {
		return errors.New("download failed")
	}

	successCh <- true
	return nil
}
