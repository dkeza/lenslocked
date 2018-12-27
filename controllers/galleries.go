package controllers

import (
	"lenslocked/models"
	"lenslocked/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// Galleries controller
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
