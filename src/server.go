package src

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anacrolix/torrent"
)

const timeoutDuration = 5 * time.Minute

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/download", downloadHandler)
	http.ListenAndServe(":8080", mux)
	log.Println("Server is listening on port 8080...")
}

func validateRequestBody(data map[string]interface{}) bool {
	_, magnet := data["magnet"].(string)
	_, folder := data["folder"].(string)
	_, name := data["name"].(string)
	_, posterUrl := data["posterUrl"].(string)
	return magnet && folder && name && posterUrl
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// verify request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse request body
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	// validate request Body
	if !validateRequestBody(data) {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// parse request body
	magnet, ok1 := data["magnet"].(string)
	folder, ok2 := data["folder"].(string)
	name, ok3 := data["name"].(string)
	posterUrl, ok4 := data["posterUrl"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		http.Error(w, "Invalid data types", http.StatusBadRequest)
		return
	}
	folderPath, pathErr := GetFolderPath(folder, name)
	if pathErr != nil {
		http.Error(w, pathErr.Error(), http.StatusBadRequest)
		return
	}

	// download torrent
	errCh := make(chan error)
	successCh := make(chan bool)
	go downloadTorrent(magnet, folderPath, errCh, successCh)

	select {
	// send any error that comes up and stop the goroutine
	case err := <-errCh:
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case <-successCh: // Capture successful completion

		// download poster
		posterErr := GetPoster(posterUrl, folderPath, name)
		if posterErr != nil {
			http.Error(w, posterErr.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Download completed successfully!"))
	// send timeout message if download takes longer than 5 minutes
	case <-time.After(timeoutDuration):
		w.Write([]byte("Download initiated or taking longer than expected"))
		return
	}
}

func handleWriteError(err error, errCh chan<- error) {
	log.Println(err)
	errCh <- err
}

func downloadTorrent(magnet string, saveDirectory string, errCh chan<- error, successCh chan<- bool) {
	// initialize torrent client
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = saveDirectory
	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		errCh <- err
		return
	}

	// add torrent
	t, err := client.AddMagnet(magnet)
	if err != nil {
		errCh <- err
		return
	}

	// wait for torrent to be ready
	<-t.GotInfo()
	log.Printf("Downloading %s...\n", t.Name())
	t.SetOnWriteChunkError(func(err error) {
		handleWriteError(err, errCh)
	})
	t.DownloadAll()

	// wait for download to complete
	ok := client.WaitAll()
	if !ok {
		errCh <- err
		return
	} else {
		client.Close()
		successCh <- true
	}
}
