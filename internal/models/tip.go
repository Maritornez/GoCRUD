package models

type Tip struct {
	Id      int    `reindex:"id,,pk"   json:"id"`
	ManId   int    `reindex:"man_id"   json:"man_id"   binding:"required"`
	Title   string `reindex:"title"    json:"title"    binding:"required"`
	Content string `reindex:"content"  json:"content"  binding:"required,max=255"`
	Pages   []struct {
		Title   string `reindex:"title"    json:"title"    binding:"required"`
		Content string `reindex:"content"  json:"content"  binding:"required,max=255"`
	} `reindex:"pages"    json:"pages"    binding:"required"`
}
