package models

// Define struct with reindex tags
type Man struct {
	ID   int64  `reindex:"id,,pk"   json:"id"`                      // 'id' is primary key
	Name string `reindex:"name"     json:"name" binding:"required"` // add index by 'name' field
	Age  int    `reindex:"age,tree" json:"age"  binding:"required"` // add sortable index by 'age' field
	Tips []Tip  `reindex:"tips"     json:"tips" binding:"required"` // add index by articles 'articles' array
}
