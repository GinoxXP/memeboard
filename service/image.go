package service

import (
	"context"
	"memeboard/models"
	"path"

	"github.com/jackc/pgx/v5"
)

const IMAGES_DIR string = "images"

func GetImagesPath() string {
	return path.Join(path.Base(""), STORAGE_DIR, IMAGES_DIR)
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

func UploadImage(imagePath string) error {
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

	err = CreateThumbnail(id)
	return err
}
