package handlers

import (
	"InvoltaTask/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"

	// use Reindexer as standalone server and connect to it via TCP.
	"github.com/restream/reindexer/v3"
	// use Reindexer as standalone server and connect to it via TCP.
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

func CreateMan(c *gin.Context) {
	// Read input
	var input models.Man
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the highest "id" value to calculate the "id" value of the added document (elem_with_highest_id.ID + 1)
	query := models.DB.Query("mans").Select("*").Sort("id", true).Limit(1)
	iterator := query.Exec()
	defer iterator.Close()

	hasNext := iterator.Next()
	var greatestId int64
	if hasNext {
		elem := iterator.Object().(*models.Man)
		greatestId = elem.ID + 1
	} else {
		greatestId = 1
	}
	// Make "man" model whitch will be added to namespace "mans"
	man := models.Man{
		ID:   greatestId,
		Name: input.Name,
		Age:  input.Age,
		Tips: input.Tips,
	}
	models.DB.Upsert("mans", man)

	// Add new document to namespace "mans"
	c.JSON(http.StatusOK, gin.H{"data": man})
}

func FindMans(c *gin.Context) {
	query := models.DB.Query("mans")
	iterator := query.Exec()
	defer iterator.Close()

	// Count amount of documents in namespace "mans"
	ns := "mans"
	iterator1 := models.DB.Query(ns).Select("*").Exec()
	defer iterator1.Close()
	amountOfDocumentsInMans := iterator.Count()

	// Create output slice "mans"
	mans := make([]models.Man, amountOfDocumentsInMans)
	// Iterate over results
	for i := 0; i < iterator.Count(); i++ {
		iterator.Next()
		// Get the next document and cast it to a pointer
		mans[i] = *iterator.Object().(*models.Man)
	}

	// Check the error
	if err := iterator.Error(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": mans})
}

func FindMan(c *gin.Context) {
	elem, found := models.DB.Query("mans").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	if found {
		item := elem.(*models.Man)
		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Item is not Found"})
	}
}

func UpdateMan(c *gin.Context) {
	elem, found := models.DB.Query("mans").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	if found {
		item := elem.(*models.Man) // founded element
		var input models.Man
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Creating model to send it to the database
		updatedMan := models.Man{
			ID:   item.ID,
			Name: input.Name,
			Age:  input.Age,
			Tips: input.Tips,
		}

		if _, err := models.DB.Update("mans", updatedMan); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Item is not Found"})
	}
}

func DeleteMan(c *gin.Context) {
	elem, found := models.DB.Query("mans").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	item := elem.(*models.Man) // founded element
	if found {
		if err := models.DB.Delete("mans", item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Item is not Found"})
	}
}
