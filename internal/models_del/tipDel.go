package models_del

type TipDel struct {
	Id    int    `reindex:"id,,pk"   json:"id"`
	Title string `reindex:"title"    json:"title"    binding:"required"`
	Pages []Page `reindex:"pages"    json:"pages"    binding:"required"`
}
type Page struct {
	Title   string `reindex:"title"    json:"title"    binding:"required"`
	Content string `reindex:"content"  json:"content"  binding:"required,max=255"`
}
