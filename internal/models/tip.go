package models

// Define struct with reindex tags
type Tip struct {
	Title   string `reindex:"title"    json:"title"    binding:"required"`
	Content string `reindex:"content"  json:"content"  binding:"required,max=255"`
	Pages   []Page `reindex:"pages"    json:"pages"    binding:"required"`
}
