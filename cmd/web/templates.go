package main

import "snippet.lciamp.xyz/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
