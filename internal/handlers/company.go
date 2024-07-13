package handlers

import (
	"net/http"
	"strconv"

	"github.com/Maritornez/GoCRUD/internal/context"
	"github.com/Maritornez/GoCRUD/internal/models"

	"github.com/gin-gonic/gin"

	// use Reindexer as standalone server and connect to it via TCP.

	// use Reindexer as standalone server and connect to it via TCP.
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

func CreateCompany(c *gin.Context) {
	// Чтение json'а, который пришел с клиента
	var input models.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Нахождение наибольшего значения id, чтобы вычислить значние id добавляемого документа
	iterator := context.DB.Query("company").Select("*").Sort("id", true).Limit(1).Exec()
	defer iterator.Close()

	hasNext := iterator.Next()
	var greatestId int
	if hasNext {
		elem := iterator.Object().(*models.Company)
		greatestId = elem.Id + 1
	} else {
		greatestId = 1
	}

	// Создание модели "company", которая будет добавлена в пространство имен "company"
	company := models.Company{
		Id:          greatestId,
		Name:        input.Name,
		Established: input.Established,
	}

	// Добавление нового документа в пространство имен "company"
	context.DB.Upsert("company", company)

	c.JSON(http.StatusOK, gin.H{"data": company})
}

func FindCompanies(c *gin.Context) {
	// Получить параметры limit и offset из запроса
	limitParam := c.DefaultQuery("limit", "10")
	offsetParam := c.DefaultQuery("offset", "0")

	// Преобразовать параметры в целые числа
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	if limit <= 0 || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit or offset parameter"})
		return
	}

	// Выполнить запрос с сортировкой по полю sort в обратном порядке
	iterator := context.DB.Query("company").Sort("name", false).Limit(limit).Offset(offset).Exec()
	defer iterator.Close()

	companies := make([]models.Company, 0)
	for iterator.Next() {
		companies = append(companies, *iterator.Object().(*models.Company))
	}

	if err := iterator.Error(); err != nil {
		panic(err)
	}

	totalCount := context.DB.Query("company").Select("*").Exec().Count()

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":        companies,
		"total_count": totalCount,
		"page":        (offset / limit) + 1,
		"total_pages": totalPages,
	})
}

func FindCompany(c *gin.Context) {
	elem, found := context.DB.Query("company").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	if found {
		item := elem.(*models.Company)
		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.JSON(http.StatusOK, gin.H{"error:": "Company with specified Id does not exist"})
	}

}

func UpdateCompany(c *gin.Context) {
	elem, found := context.DB.Query("company").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	if found {
		item := elem.(*models.Company)
		var input struct {
			Name        *string `json:"name"`
			Established *int    `json:"established"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Обновление только тех полей, которые были переданы в json
		if input.Name != nil {
			item.Name = *input.Name
		}
		if input.Established != nil {
			item.Established = *input.Established
		}

		if _, err := context.DB.Update("company", item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Company with specified Id does not exist"})
	}
}

func DeleteCompany(c *gin.Context) {
	elem, found := context.DB.Query("company").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	if found {
		item := elem.(*models.Company)
		if err := context.DB.Delete("company", item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Company with specified Id does not exist"})
	}
}
