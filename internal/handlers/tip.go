package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/restream/reindexer/v3" // use Reindexer as standalone server and connect to it via TCP.
	_ "github.com/restream/reindexer/v3/bindings/cproto"

	"github.com/Maritornez/GoCRUD/internal/models"
	"github.com/Maritornez/GoCRUD/internal/models_del"
	"github.com/Maritornez/GoCRUD/internal/storage"
)

func CreateTip(c *gin.Context) {
	var input models.Tip
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Проверка наличия man с указанным id
	_, found := storage.DB.Query("man").
		Where("id", reindexer.EQ, input.ManId).
		Get()
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Man with specified Id does not exist"})
		return
	}

	// Нахождение наибольшего значения id, чтобы вычислить значние id добавляемого документа
	iterator := storage.DB.Query("tip").Select("*").Sort("id", true).Limit(1).Exec()
	defer iterator.Close()

	hasNext := iterator.Next()
	var greatestId int
	if hasNext {
		elem := iterator.Object().(*models.Tip)
		greatestId = elem.Id + 1
	} else {
		greatestId = 1
	}

	// Создание модели "tip", которая будет добавлена в пространство имен "tip"
	tip := models.Tip{
		Id:    greatestId,
		ManId: input.ManId,
		Title: input.Title,
		Pages: input.Pages,
	}

	// Добавление нового документа в пространство имен "man"
	storage.DB.Upsert("tip", tip)

	c.JSON(http.StatusOK, gin.H{"data": tip})
}

// Функция для обработки модели Tip (удаляется поле manId)
func processTip(tip models.Tip) models_del.TipDel {
	pagesForTipDel := make([]models_del.Page, 0)

	for _, page := range tip.Pages {
		pagesForTipDel = append(pagesForTipDel, models_del.Page{
			Title:   page.Title,
			Content: page.Content,
		})
	}

	tipDel := models_del.TipDel{
		Id:    tip.Id,
		Title: tip.Title,
		Pages: pagesForTipDel,
	}
	return tipDel
}

func FindTips(c *gin.Context) {
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

	// Выполнить запрос с сортировкой по полю man_id в обратном порядке
	iterator := storage.DB.Query("tip").Sort("man_id", false).Limit(limit).Offset(offset).Exec()
	defer iterator.Close()

	var wg sync.WaitGroup
	tipsDel := make([]models_del.TipDel, 0)
	tipDelChannel := make(chan models_del.TipDel, limit)
	for iterator.Next() {
		tip := *iterator.Object().(*models.Tip)
		wg.Add(1)

		//tips = append(tips, *iterator.Object().(*models.Tip))

		// Обработка каждого документа в отдельной горутине
		go func(tip models.Tip) {
			defer wg.Done()
			tipDelChannel <- processTip(tip)
		}(tip)
	}

	go func() {
		wg.Wait()
		close(tipDelChannel)
	}()

	for tipDel := range tipDelChannel {
		tipsDel = append(tipsDel, tipDel)
	}

	if err := iterator.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalCount := storage.DB.Query("tip").Select("*").Exec().Count()

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":        tipsDel,
		"total_count": totalCount,
		"page":        (offset / limit) + 1,
		"total_pages": totalPages,
	})
}

func FindTip(c *gin.Context) {
	//Проверка наличия tip с указанным id в кэше
	id := c.Param("id")
	cacheKey := "tip:" + id

	if cachedData, err := storage.Cache.Get(cacheKey); err == nil {
		var tip models.Tip
		if err := json.Unmarshal(cachedData, &tip); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": tip})
			return
		}
	}

	//Проверка наличия tip с указанным id в базе данных
	elem, found := storage.DB.Query("tip").Where("id", reindexer.EQ, c.Param("id")).Get()

	if found {
		item := elem.(*models.Tip)

		// Кэширование данных
		jsonData, err := json.Marshal(item)
		if err == nil {
			storage.Cache.Set(cacheKey, jsonData)
		}

		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Tip with specified Id does not exist"})
		return
	}
}

func UpdateTip(c *gin.Context) {
	elem, found := storage.DB.Query("tip").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	if found {
		item := elem.(*models.Tip)
		var input struct {
			ManId *int    `json:"man_id"`
			Title *string `json:"title"`
			Pages *[]struct {
				Title   *string `json:"title"`
				Content *string `json:"content"`
			} `json:"pages"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Обновление только тех полей, которые были переданы в json
		if input.ManId != nil {
			_, found := storage.DB.Query("man").
				Where("id", reindexer.EQ, input.ManId).
				Get()
			if !found {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Man with specified Id does not exist"})
				return
			}

			item.ManId = *input.ManId
		}
		if input.Title != nil {
			item.Title = *input.Title
		}
		if input.Pages != nil {
			item.Pages = make([]models.Page, len(*input.Pages))
			for i, inputPage := range *input.Pages {
				if inputPage.Title != nil {
					item.Pages[i].Title = *inputPage.Title
				}
				if inputPage.Content != nil {
					item.Pages[i].Content = *inputPage.Content
				}
			}
		}

		if _, err := storage.DB.Update("tip", item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Tip with specified Id does not exist"})
		return
	}
}

func DeleteTip(c *gin.Context) {
	elem, found := storage.DB.Query("tip").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()
	if found {
		item := elem.(*models.Tip)
		if err := storage.DB.Delete("tip", item); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Tip with specified Id does not exist"})
		return
	}
}
