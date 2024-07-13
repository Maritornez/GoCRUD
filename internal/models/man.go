package models

type Man struct {
	Id        int    `reindex:"id,,pk" json:"id"`
	Name      string `reindex:"name" json:"name" binding:"required"`    // add index by 'name' field
	Age       int    `reindex:"age,tree" json:"age" binding:"required"` // add sortable index by 'age' field
	CompanyId int    `reindex:"company_id" json:"company_id" binding:"required"`
	Sort      int    `reindex:"sort" json:"sort" binding:"required"`
	//Tips []Tip `reindex:"tips" json:"tips" binding:"required"`   // add index by articles 'articles' array
}
