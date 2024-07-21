package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/Maritornez/GoCRUD/internal/models"
	"github.com/Maritornez/GoCRUD/internal/models_bind"
	"github.com/Maritornez/GoCRUD/internal/models_del"
	"github.com/Maritornez/GoCRUD/internal/storage"

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

	//Проверка наличия company с указанным id
	_, found := storage.DB.Query("company").
		Where("id", reindexer.EQ, input.CompanyId).
		Get()
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Company with specified Id does not exist"})
		return
	}

	// Нахождение наибольшего значения id, чтобы вычислить значние id добавляемого документа
	iterator := storage.DB.Query("man").Select("*").Sort("id", true).Limit(1).Exec()
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

	storage.DB.Upsert("man", man)

	c.JSON(http.StatusOK, gin.H{"data": man})
}

// Функция для обработки модели Man (удаляется поля companyId)
// (вложенные tip тоже обрабатываются)
func processMan(man models_bind.ManBind) models_del.ManBindDel {
	processedTips := make([]models_del.TipDel, 0)
	for i := range man.Tips {
		processedTips = append(processedTips, processTip(man.Tips[i]))
	}

	manDel := models_del.ManBindDel{
		Id:   man.Id,
		Name: man.Name,
		Age:  man.Age,
		Sort: man.Sort,
		Tips: processedTips,
	}
	return manDel
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
	iterator := storage.DB.Query("man").Sort("sort", true).Limit(limit).Offset(offset).Exec()
	defer iterator.Close()

	var wg sync.WaitGroup
	menBindDel := make([]models_del.ManBindDel, 0)
	manBindDelChannel := make(chan models_del.ManBindDel, limit)
	for iterator.Next() {
		man := iterator.Object().(*models.Man)
		wg.Add(1)

		// Поиск подходящих tip для данного man
		iterator := storage.DB.Query("tip").Sort("title", false).Where("man_id", reindexer.EQ, man.Id).Exec()
		defer iterator.Close()

		tips := make([]models.Tip, 0)
		for iterator.Next() {
			tips = append(tips, *iterator.Object().(*models.Tip))
		}

		manBind := models_bind.ManBind{
			Id:        man.Id,
			Name:      man.Name,
			Age:       man.Age,
			CompanyId: man.CompanyId,
			Sort:      man.Sort,
			Tips:      tips,
		}
		//menBind = append(menBind, manBind)

		// Обработка каждого документа в отдельной горутине
		go func(man models_bind.ManBind) {
			defer wg.Done()
			manBindDelChannel <- processMan(manBind)
		}(manBind)
	}

	go func() {
		wg.Wait()
		close(manBindDelChannel)
	}()

	for manBindDel := range manBindDelChannel {
		menBindDel = append(menBindDel, manBindDel)
	}

	if err := iterator.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalCount := storage.DB.Query("man").Select("*").Exec().Count()

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":        menBindDel,
		"total_count": totalCount,
		"page":        (offset / limit) + 1,
		"total_pages": totalPages,
	})
}

func FindMan(c *gin.Context) {
	//Проверка наличия tip с указанным id в кэше
	id := c.Param("id")
	cacheKey := "man:" + id

	if cachedData, err := storage.Cache.Get(cacheKey); err == nil {
		var manBind models_bind.ManBind
		if err := json.Unmarshal(cachedData, &manBind); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": manBind})
			return
		}
	}

	// elem is interface{}
	elem, found := storage.DB.Query("man").Where("id", reindexer.EQ, c.Param("id")).Get()

	if found {
		item := elem.(*models.Man) // item is *models.Man

		// Поиск подходящих tip для данного man
		iterator := storage.DB.Query("tip").Where("man_id", reindexer.EQ, item.Id).Sort("title", false).Exec()
		defer iterator.Close()

		tips := make([]models.Tip, 0)
		for iterator.Next() {
			tips = append(tips, *iterator.Object().(*models.Tip))
		}

		if err := iterator.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		manBind := models_bind.ManBind{
			Id:        item.Id,
			Name:      item.Name,
			Age:       item.Age,
			CompanyId: item.CompanyId,
			Sort:      item.Sort,
			Tips:      tips,
		}

		// Кэширование данных
		jsonData, err := json.Marshal(manBind)
		if err == nil {
			storage.Cache.Set(cacheKey, jsonData)
		}

		c.JSON(http.StatusOK, gin.H{"data": manBind})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Man with specified Id does not exist"})
		return
	}
}

func UpdateMan(c *gin.Context) {
	elem, found := storage.DB.Query("man").
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
			_, found := storage.DB.Query("company").
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

		if _, err := storage.DB.Update("man", item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Man with specified Id does not exist"})
		return
	}
}

func DeleteMan(c *gin.Context) {
	elem, found := storage.DB.Query("man").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()
	if found {
		item := elem.(*models.Man)

		// Удаление связанных tip
		iterator := storage.DB.Query("tip").Where("man_id", reindexer.EQ, item.Id).Exec()
		defer iterator.Close()

		for iterator.Next() {
			if err := storage.DB.Delete("tip", iterator.Object().(*models.Tip)); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		if err := storage.DB.Delete("man", item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Man with specified Id does not exist"})
		return
	}
}
