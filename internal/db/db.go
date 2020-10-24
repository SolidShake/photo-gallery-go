package db

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	DB, err = sqlx.Open("sqlite3", "./internal/db/data/gallery.db")
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS `images` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `filename` VARCHAR(255) NOT NULL, `extension` VARCHAR(20) NOT NULL, `mimetype` VARCHAR(20) NOT NULL)")
	if err != nil {
		panic(err)
	}
}

func SaveImage(filename, extension, mimetype string) (err error) {
	_, err = DB.Exec("INSERT INTO images (filename, extension, mimetype) VALUES ($1, $2, $3)", filename, extension, mimetype)
	if err != nil {
		log.Println(err)
	}
	return
}

func DeleteImage(filename string) (err error) {
	res, err := DB.Exec("DELETE FROM images WHERE filename = $1", filename)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == int64(0) {
				return errors.New("Images not found")
			}
		}
	}
	return
}

func GetImageByHash(hash string) (image Image, err error) {
	err = DB.Get(&image, "SELECT * FROM images where filename = $1", hash)
	if err != nil {
		log.Println(err)
	}
	return
}

func GetAllImages() (images []Image, err error) {
	err = DB.Select(&images, "SELECT * FROM images")
	if err != nil {
		log.Println(err)
	}
	return
}
