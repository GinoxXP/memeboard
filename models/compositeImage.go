package models

type CompositeImage struct {
	Image     *Image
	Thumbnail *Thumbnail
}
