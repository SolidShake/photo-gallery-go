package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SolidShake/photo-gallery-go/internal/db"
	"github.com/SolidShake/photo-gallery-go/internal/handler"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	db.InitDB()
	defer db.DB.Close()

	router := httprouter.New()
	router.GET("/photos", handler.GetPhotoList)
	router.POST("/photos/upload", handler.UploadFile)
	router.DELETE("/photos/:hash", handler.DeletePhoto)
	router.GET("/photos/:hash", handler.GetPhoto)
	router.GET("/photos/:hash/preview", handler.GetPhotoPreview)

	log.Println("Service started at port: " + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), router))
}

func initConfig() error {
	workingdir, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.SetConfigFile(workingdir + "/configs/config.yml")
	return viper.ReadInConfig()
}
