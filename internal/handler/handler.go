package handler

import (
	"crypto/rand"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SolidShake/photo-gallery-go/internal/db"
	"github.com/SolidShake/photo-gallery-go/internal/response"
	"github.com/disintegration/imaging"
	"github.com/julienschmidt/httprouter"
)

const maxUploadSize = 2 * 1024 * 1024 // 2 mb
const uploadPath = "./storage"

const previewHeight = 500 // 500px

func UploadFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Println(err)
		response.RenderResponse(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("uploadFile")
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileSize := fileHeader.Size
	if fileSize > maxUploadSize {
		response.RenderResponse(w, "File too big", http.StatusBadRequest)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "Invalid file", http.StatusBadRequest)
		return
	}

	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg":
	case "image/jpg":
	case "image/png":
	case "image/gif":
		break
	default:
		response.RenderResponse(w, "Invalid file MIME type", http.StatusBadRequest)
		return
	}
	fileName := randToken(12)
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "Can't parse file extension", http.StatusBadRequest)
		return
	}
	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])

	newFile, err := os.Create(newPath)
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "Can't save file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		response.RenderResponse(w, "Can't save file", http.StatusInternalServerError)
		return
	}

	err = db.SaveImage(fileName, fileEndings[0], detectedFileType)
	if err != nil {
		response.RenderResponse(w, "Can't save file", http.StatusInternalServerError)
		return
	}

	response.RenderPhotoUrlResponse(w, r, fileName, http.StatusCreated)
}

func DeletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dbImg, err := getPhotoInfo(ps.ByName("hash"))
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "Photo not found", http.StatusBadRequest)
		return
	}

	err = db.DeleteImage(ps.ByName("hash"))
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "Photo not found", http.StatusBadRequest)
		return
	}

	err = os.Remove(getPhotoPath(dbImg))
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "Photo not found", http.StatusBadRequest)
	}

	response.RenderResponse(w, "Photo deleted", http.StatusOK)
}

func GetPhotoList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	images, err := db.GetAllImages()
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.RenderPhotosResponse(w, r, images, http.StatusOK)
}

func GetPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	img, mime, err := getPhoto(ps.ByName("hash"))
	defer img.Close()
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "photo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mime)
	io.Copy(w, img)
}

func GetPhotoPreview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	preview, dbImg, err := getPhotoPreview(ps.ByName("hash"))
	if err != nil {
		log.Println(err)
		response.RenderResponse(w, "photo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", dbImg.Mimetype)
	format, _ := imaging.FormatFromExtension(dbImg.Extension)
	imaging.Encode(w, preview, format)
}

func getPhotoPreview(hash string) (image *image.NRGBA, dbImg db.Image, err error) {
	dbImg, err = getPhotoInfo(hash)
	if err != nil {
		return
	}

	img, err := imaging.Open(getPhotoPath(dbImg), imaging.AutoOrientation(true))
	image = imaging.Resize(img, previewHeight, 0, imaging.Lanczos)
	return
}

func getPhoto(hash string) (img *os.File, mime string, err error) {
	dbImg, err := getPhotoInfo(hash)
	if err != nil {
		return
	}
	mime = dbImg.Mimetype

	img, err = os.Open(getPhotoPath(dbImg))
	if err != nil {
		return
	}

	return
}

func getPhotoInfo(hash string) (dbImg db.Image, err error) {
	return db.GetImageByHash(hash)
}

func getPhotoPath(dbImg db.Image) string {
	return uploadPath + "/" + dbImg.Filename + dbImg.Extension
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
