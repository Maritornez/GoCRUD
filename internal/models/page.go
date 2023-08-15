package models

type Page struct {
	Title   string `reindex:"title"    json:"title"    binding:"required"`
	Content string `reindex:"content"  json:"content"  binding:"required,max=255"`
}
