package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nfnt/resize"
)

const STORAGE_DIR string = "storage"
const IMAGES_DIR string = "images"
const THUMBNAILS_DIR string = "thumbnails"
const DB_SETTINGS_FILE string = "dbSettings.json"

type Storage struct {
	connection pgxpool.Pool
}

func NewStorage() *Storage {
	return &Storage{
		connection: *ConnectToDB(),
	}
}

func ConnectToDB() *pgxpool.Pool {
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

	conn, err := pgxpool.New(context.Background(), bdUrl)
	if err != nil {
		panic(err)
	}

	return conn
}

func GetImagesPath() string {
	return path.Join(path.Base(""), STORAGE_DIR, IMAGES_DIR)
}

func GetThumbnailsPath() string {
	return path.Join(path.Base(""), STORAGE_DIR, THUMBNAILS_DIR)
}

func GetDbSettings() string {
	return path.Join(path.Base(""), DB_SETTINGS_FILE)
}

func (storage *Storage) GetAllImages() ([]string, error) {
	connection := &storage.connection

	rows, err := connection.Query(
		context.Background(),
		"SELECT path FROM image")
	if err != nil {
		return nil, err
	}

	paths, err := pgx.CollectRows(rows, pgx.RowTo[string])

	return paths, err
}

func (storage *Storage) GetAllThumbnails() ([]string, error) {
	connection := &storage.connection

	rows, err := connection.Query(
		context.Background(),
		"SELECT path FROM thumbnail")
	if err != nil {
		return nil, err
	}

	paths, err := pgx.CollectRows(rows, pgx.RowTo[string])

	return paths, err
}

func (storage *Storage) UploadImage(imagePath string) error {
	connection := &storage.connection

	idRows, err := connection.Query(
		context.Background(),
		"INSERT INTO image (path) VALUES ($1) RETURNING id;", imagePath)
	if err != nil {
		return err
	}

	id, err := pgx.CollectOneRow(idRows, pgx.RowTo[int])
	if err != nil {
		return err
	}

	err = storage.CreateThumbnail(id)
	return err
}

func (storage *Storage) CreateThumbnail(imageId int) error {
	connection := &storage.connection

	pathRows, err := connection.Query(
		context.Background(),
		"SELECT path FROM image WHERE id = $1;", imageId)
	if err != nil {
		return err
	}

	imagePath, err := pgx.CollectOneRow(pathRows, pgx.RowTo[string])
	if err != nil {
		return err
	}

	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sourceImage, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	thumbnail := resize.Thumbnail(250, 250, sourceImage, resize.Bilinear)

	thumbnailName := uuid.New().String() + ".png"
	thumbnailPath := path.Join(GetThumbnailsPath(), thumbnailName)

	thumbnailFile, err := os.Create(thumbnailPath)
	if err != nil {
		return err
	}
	defer thumbnailFile.Close()

	err = png.Encode(thumbnailFile, thumbnail)
	if err != nil {
		return err
	}

	thumbnailIdRows, err := connection.Query(
		context.Background(),
		"INSERT INTO thumbnail (path) VALUES ($1) RETURNING id;", thumbnailPath)
	if err != nil {
		return err
	}

	thumbnailId, err := pgx.CollectOneRow(thumbnailIdRows, pgx.RowTo[int])
	if err != nil {
		return err
	}

	_, err = connection.Query(
		context.Background(),
		"INSERT INTO image_thumbnail (image_id, thumbnail_id) VALUES ($1, $2);", imageId, thumbnailId)
	return err
}
