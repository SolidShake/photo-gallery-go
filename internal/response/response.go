package response

import (
	"encoding/json"
	"net/http"

	"github.com/SolidShake/photo-gallery-go/internal/db"
)

type ResponseErrorBody struct {
	Message string
}

type ResponseUrlBody struct {
	PhotoUrl        string
	PhotoPreviewUrl string
}

func RenderResponse(w http.ResponseWriter, message string, statusCode int) {
	responseBody := ResponseErrorBody{Message: message}
	data, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func RenderPhotoUrlResponse(w http.ResponseWriter, r *http.Request, hash string, statusCode int) {
	photoURL := r.Host + "/photos/" + hash
	photoPreviewURL := r.Host + "/photos/" + hash + "/preview"
	responseBody := ResponseUrlBody{PhotoUrl: photoURL, PhotoPreviewUrl: photoPreviewURL}
	data, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func RenderPhotosResponse(w http.ResponseWriter, r *http.Request, images []db.Image, statusCode int) {
	imagesResponseBody := make([]ResponseUrlBody, len(images))
	for i, img := range images {
		imagesResponseBody[i].PhotoUrl = r.Host + "/photos/" + img.Filename
		imagesResponseBody[i].PhotoPreviewUrl = r.Host + "/photos/" + img.Filename + "/preview"
	}
	data, err := json.Marshal(imagesResponseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
