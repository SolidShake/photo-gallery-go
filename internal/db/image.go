package db

type Image struct {
	ID        int    `db:"id"`
	Filename  string `db:"filename"`
	Extension string `db:"extension"`
	Mimetype  string `db:"mimetype"`
}
