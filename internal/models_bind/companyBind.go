package models_bind

type CompanyBind struct {
	Id          int    `reindex:"id,,pk" json:"id"`
	Name        string `reindex:"name" json:"name" binding:"required"`
	Established int    `reindex:"established" json:"established" binding:"required"`
	// Связанное поле
	Men []ManBind `reindex:"men" json:"men"`
}
