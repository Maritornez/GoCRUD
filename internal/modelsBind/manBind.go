package modelsBind

import (
	"github.com/Maritornez/GoCRUD/internal/models"
)

type ManBind struct {
	Id        int    `reindex:"id,,pk" json:"id"`
	Name      string `reindex:"name" json:"name" binding:"required"`
	Age       int    `reindex:"age,tree" json:"age" binding:"required"`
	CompanyId int    `reindex:"company_id" json:"company_id" binding:"required"`
	Sort      int    `reindex:"sort" json:"sort" binding:"required"`
	// Связанное поле
	Tips []models.Tip `reindex:"tips" json:"tips"`
}
