package service

import (
	"context"
	"errors"
	"fmt"
	"memeboard/models"

	"github.com/jackc/pgx/v5"
)

func GetAllTags() ([]string, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM tag")
	if err != nil {
		return nil, err
	}

	tags, err := pgx.CollectRows(rows, pgx.RowTo[string])

	return tags, err
}

func GetAttachedTags(imageID int) ([]*models.Tag, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT tag_id FROM image_tag WHERE image_id = $1", imageID)
	if err != nil {
		return nil, err
	}

	tagIDs, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return nil, err
	}

	tags := make([]*models.Tag, len(tagIDs))
	for i, id := range tagIDs {
		tag, err := GetTagByID(id)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}
	return tags, err
}

func AddTag(tagname string) (*models.Tag, error) {
	rows, err := connection.Query(
		context.Background(),
		"INSERT INTO tag (tagname) VALUES ($1) RETURNING *", tagname)
	if err != nil {
		return nil, err
	}

	tag, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Tag])
	return &tag, err
}

func AddTagToImage(tagname string, imageID int) error {
	isTagExists, err := TagExists(tagname)

	if err != nil {
		return err
	}

	var tag *models.Tag

	if isTagExists {
		tag, err = GetTagByName(tagname)
	} else {
		tag, err = AddTag(tagname)
	}

	if err != nil {
		return err
	}

	isImageTagLinkExist, err := TagImageLinkExists(imageID, tag.ID)
	if err != nil {
		return err
	}

	if isImageTagLinkExist {
		message := fmt.Sprintf("Tag %s already attached to image with ID %d", tagname, imageID)
		return errors.New(message)
	}

	err = CreateTagImageLink(imageID, tag.ID)
	return err
}

func RemoveTagFromImage(tagId int, imageId int) error {
	_, err := connection.Query(
		context.Background(),
		"DELETE FROM image_tag WHERE tag_id = $1 AND image_id = $2", tagId, imageId)

	return err
}

func TagImageLinkExists(imageId int, tagId int) (bool, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT EXISTS(SELECT * FROM image_tag WHERE image_id = $1 AND tag_id = $2)", imageId, tagId)
	if err != nil {
		return false, err
	}

	isExist, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	return isExist, err
}

func CreateTagImageLink(imageId int, tagId int) error {
	_, err := connection.Query(
		context.Background(),
		"INSERT INTO image_tag (image_id, tag_id) VALUES ($1, $2)", imageId, tagId)

	return err
}

func TagExists(tag string) (bool, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT EXISTS(SELECT * FROM tag WHERE tagname = ($1))", tag)

	if err != nil {
		return false, err
	}

	isExist, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])

	return isExist, err
}

func GetTagByName(tagname string) (*models.Tag, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM tag WHERE tagname = '($1)'", tagname)

	if err != nil {
		return nil, err
	}

	tag, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Tag])
	return &tag, err
}

func GetTagByID(id int) (*models.Tag, error) {
	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM tag WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	tag, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Tag])
	return &tag, err
}
