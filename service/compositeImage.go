package service

import (
	"memeboard/models"
)

func GetCompositeImages(images *[]models.Image) (*[]models.CompositeImage, error) {
	compositeImages := make([]models.CompositeImage, len(*images))

	for i, image := range *images {
		thumnnail, err := GetThumbnail(image.ID)
		if err != nil {
			return nil, err
		}

		compositeImages[i] = models.CompositeImage{
			Image:     &image,
			Thumbnail: thumnnail,
		}
	}

	return &compositeImages, nil
}
