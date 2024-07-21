package models_del

type ManBindDel struct {
	Id   int    `reindex:"id,,pk" json:"id"`
	Name string `reindex:"name" json:"name" binding:"required"`
	Age  int    `reindex:"age,tree" json:"age" binding:"required"`
	Sort int    `reindex:"sort" json:"sort" binding:"required"`
	// Связанное поле
	Tips []TipDel `reindex:"tips" json:"tips"`
}
