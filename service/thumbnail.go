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

const THUMBNAILS_DIR string = "thumbnails"

func GetThumbnailsPath() string {
	return path.Join(path.Base(""), STORAGE_DIR, THUMBNAILS_DIR)
}

func GetThumbnail(imageId int) (*models.Thumbnail, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM thumbnail "+
			"WHERE id = ( "+
			"SELECT thumbnail_id FROM image_thumbnail "+
			"WHERE image_id = $1)", imageId)
	if err != nil {
		return nil, err
	}

	thumbnails, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Thumbnail])

	return &thumbnails[0], err
}

func CreateThumbnail(imageId int) error {
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
