package models

type Man struct {
	Id        int    `reindex:"id,,pk" json:"id"`
	Name      string `reindex:"name" json:"name" binding:"required"`
	Age       int    `reindex:"age,tree" json:"age" binding:"required"`
	CompanyId int    `reindex:"company_id" json:"company_id" binding:"required"`
	Sort      int    `reindex:"sort" json:"sort" binding:"required"`
}
