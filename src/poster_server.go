package src

import (
	"encoding/json"
	"log"
	"net/http"
)

func setPosterHandler(w http.ResponseWriter, r *http.Request) {
	// validate request body
	var req SetPosterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, err, 400)
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, err, 400)
		return
	}

	log.Println("Updating poster for download with name", req.Name)

	// get folder path
	folderPath, err := GetFolderPath(req.Folder, req.Name)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	// download the poster
	if err := GetPoster(req.URL, folderPath); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
		return
	}
}
