package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/jackc/pgx/v5"
)

const STORAGE_DIR string = "storage"
const IMAGES_DIR string = "images"
const DB_SETTINGS_FILE string = "dbSettings.json"

type Storage struct {
	connection pgx.Conn
}

func NewStorage() *Storage {
	return &Storage{
		connection: ConnectToDB(),
	}
}

func ConnectToDB() pgx.Conn {
	jsonFile, err := os.Open(GetDbSettings())
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var settings DbSettings

	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		panic(err)
	}

	bdUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		settings.Username,
		settings.Password,
		settings.Address,
		settings.Port,
		settings.DbName)

	conn, err := pgx.Connect(context.Background(), bdUrl)
	if err != nil {
		panic(err)
	}

	return *conn
}

func GetImagesPath() string {
	return path.Join(path.Base(""), STORAGE_DIR, IMAGES_DIR)
}

func GetDbSettings() string {
	return path.Join(path.Base(""), DB_SETTINGS_FILE)
}

func (storage Storage) GetAllImages() ([]string, error) {
	connection := &storage.connection

	rows, err := connection.Query(context.Background(), "SELECT path FROM image")

	if err != nil {
		return nil, err
	}

	paths, err := pgx.CollectRows(rows, pgx.RowTo[string])

	return paths, err
}

func (storage Storage) UploadImage(imagePath string) error {
	connection := &storage.connection

	err := connection.QueryRow(
		context.Background(),
		"INSERT INTO image (path) VALUES ('"+imagePath+"')").Scan()

	return err
}
