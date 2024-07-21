package models_del

type CompanyBindDel struct {
	Id   int    `reindex:"id,,pk" json:"id"`
	Name string `reindex:"name" json:"name" binding:"required"`
	// Связанное поле
	Men []ManBindDel `reindex:"men" json:"men"`
}
