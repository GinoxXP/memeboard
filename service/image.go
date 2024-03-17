package service

import (
	"context"
	"image"
	"image/png"
	"memeboard/models"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nfnt/resize"
)

const IMAGES_DIR string = "images"

const THUMBNAILS_DIR string = "thumbnails"

func GetImagesPath() string {
	return path.Join(path.Base(""), STORAGE_DIR, IMAGES_DIR)
}

func GetThumbnailsPath() string {
	return path.Join(path.Base(""), STORAGE_DIR, THUMBNAILS_DIR)
}

func GetAllImages() (*[]models.Image, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM image")
	if err != nil {
		return nil, err
	}

	images, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Image])

	return &images, err
}

func GetImage(id int) (*models.Image, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM image WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	image, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Image])

	return &image, err
}

func UploadImage(imagePath string) error {
	thumbnailPath, err := CreateThumbnail(imagePath)
	_, err = connection.Query(
		context.Background(),
		"INSERT INTO image (source_path, thumbnail_path) VALUES ($1, $2)", imagePath, thumbnailPath)
	if err != nil {
		return err
	}

	return err
}

func CreateThumbnail(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sourceImage, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	thumbnail := resize.Thumbnail(250, 250, sourceImage, resize.Bilinear)

	thumbnailName := uuid.New().String() + ".png"
	thumbnailPath := path.Join(GetThumbnailsPath(), thumbnailName)

	thumbnailFile, err := os.Create(thumbnailPath)
	if err != nil {
		return "", err
	}
	defer thumbnailFile.Close()

	err = png.Encode(thumbnailFile, thumbnail)
	if err != nil {
		return "", err
	}

	return thumbnailPath, err
}
