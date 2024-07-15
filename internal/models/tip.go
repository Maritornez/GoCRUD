package models

type Tip struct {
	Id    int    `reindex:"id,,pk"   json:"id"`
	ManId int    `reindex:"man_id"   json:"man_id"   binding:"required"`
	Title string `reindex:"title"    json:"title"    binding:"required"`
	Pages []Page `reindex:"pages"    json:"pages"    binding:"required"`
}
type Page struct {
	Title   string `reindex:"title"    json:"title"    binding:"required"`
	Content string `reindex:"content"  json:"content"  binding:"required,max=255"`
}
