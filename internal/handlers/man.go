package handlers

import (
	"net/http"
	"strconv"

	"github.com/Maritornez/GoCRUD/internal/context"
	"github.com/Maritornez/GoCRUD/internal/models"

	"github.com/gin-gonic/gin"

	// use Reindexer as standalone server and connect to it via TCP.
	"github.com/restream/reindexer/v3"
	// use Reindexer as standalone server and connect to it via TCP.
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

func CreateMan(c *gin.Context) {
	// Чтение json'а, который пришел с клиента
	var input models.Man
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Проверка наличия компании с указанным id
	_, found := context.DB.Query("company").
		Where("id", reindexer.EQ, input.CompanyId).
		Get()
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Company with specified Id does not exist"})
		return
	}

	// Нахождение наибольшего значения id, чтобы вычислить значние id добавляемого документа
	iterator := context.DB.Query("man").Select("*").Sort("id", true).Limit(1).Exec()
	defer iterator.Close()

	hasNext := iterator.Next()
	var greatestId int
	if hasNext {
		elem := iterator.Object().(*models.Man)
		greatestId = elem.Id + 1
	} else {
		greatestId = 1
	}

	// Создание модели "man", которая будет добавлена в пространство имен "man"
	man := models.Man{
		Id:        greatestId,
		Name:      input.Name,
		Age:       input.Age,
		CompanyId: input.CompanyId,
		Sort:      input.Sort,
	}

	// Добавление нового документа в пространство имен "man"
	context.DB.Upsert("man", man)

	c.JSON(http.StatusOK, gin.H{"data": man})
}

func FindMen(c *gin.Context) {
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
	iterator := context.DB.Query("man").Sort("sort", true).Limit(limit).Offset(offset).Exec()
	defer iterator.Close()

	men := make([]models.Man, 0)
	for iterator.Next() {
		men = append(men, *iterator.Object().(*models.Man))
	}

	if err := iterator.Error(); err != nil {
		panic(err)
	}

	totalCount := context.DB.Query("man").Select("*").Exec().Count()

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":        men,
		"total_count": totalCount,
		"page":        (offset / limit) + 1,
		"total_pages": totalPages,
	})
}

func FindMan(c *gin.Context) {
	elem, found := context.DB.Query("man"). // elem is interface{}
						Where("id", reindexer.EQ, c.Param("id")).
						Get()

	if found {
		item := elem.(*models.Man) // item is *models.Man
		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Man with specified Id does not exist"})
	}
}

func UpdateMan(c *gin.Context) {
	elem, found := context.DB.Query("man").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	if found {
		item := elem.(*models.Man)
		var input struct {
			Name      *string `json:"name"`
			Age       *int    `json:"age"`
			CompanyId *int    `json:"company_id"`
			Sort      *int    `json:"sort"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Обновление только тех полей, которые были переданы в json
		if input.Name != nil {
			item.Name = *input.Name
		}
		if input.Age != nil {
			item.Age = *input.Age
		}
		if input.CompanyId != nil {
			_, found := context.DB.Query("company").
				Where("id", reindexer.EQ, input.CompanyId).
				Get()
			if !found {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Company with specified Id does not exist"})
				return
			}
			item.CompanyId = *input.CompanyId
		}
		if input.Sort != nil {
			item.Sort = *input.Sort
		}

		if _, err := context.DB.Update("man", item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Man with specified Id does not exist"})
	}
}

func DeleteMan(c *gin.Context) {
	elem, found := context.DB.Query("man").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()
	if found {
		item := elem.(*models.Man)
		if err := context.DB.Delete("man", item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Man with specified Id does not exist"})
	}
}
