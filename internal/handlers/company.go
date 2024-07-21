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
	"github.com/restream/reindexer/v3" // use Reindexer as standalone server and connect to it via TCP.
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
	iterator := storage.DB.Query("company").Select("*").Sort("id", true).Limit(1).Exec()
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
	storage.DB.Upsert("company", company)

	c.JSON(http.StatusOK, gin.H{"data": company})
}

// Функция для обработки модели Company (удаляется поле Established)
// (вложенные man и tip тоже обрабатываются)
func processCompany(company models_bind.CompanyBind) models_del.CompanyBindDel {
	processedMen := make([]models_del.ManBindDel, 0)
	for i := range company.Men {
		processedMen = append(processedMen, processMan(company.Men[i]))
	}

	companyBindDel := models_del.CompanyBindDel{
		Id:   company.Id,
		Name: company.Name,
		Men:  processedMen,
	}
	return companyBindDel
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

	// Выполнить запрос с сортировкой по полю name в обратном порядке
	iterator := storage.DB.Query("company").Sort("name", false).Limit(limit).Offset(offset).Exec()
	defer iterator.Close()

	var wg sync.WaitGroup
	companiesBindDel := make([]models_del.CompanyBindDel, 0)
	companyBindDelChannel := make(chan models_del.CompanyBindDel, limit)
	for iterator.Next() {
		company := iterator.Object().(*models.Company)
		wg.Add(1)

		// Поиск подходящих man для данной company
		iterator := storage.DB.Query("man").Where("company_id", reindexer.EQ, company.Id).Sort("sort", true).Exec()
		defer iterator.Close()

		menBind := make([]models_bind.ManBind, 0)
		for iterator.Next() {
			man := iterator.Object().(*models.Man)

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
			menBind = append(menBind, manBind)
		}

		if err := iterator.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		companyBind := models_bind.CompanyBind{
			Id:          company.Id,
			Name:        company.Name,
			Established: company.Established,
			Men:         menBind,
		}

		//companiesBind = append(companiesBind, companyBind)

		// Обработка каждого документа в отдельной горутине
		go func(company models_bind.CompanyBind) {
			defer wg.Done()
			companyBindDelChannel <- processCompany(company)
		}(companyBind)
	}

	go func() {
		wg.Wait()
		close(companyBindDelChannel)
	}()

	for companyBindDel := range companyBindDelChannel {
		companiesBindDel = append(companiesBindDel, companyBindDel)
	}

	if err := iterator.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalCount := storage.DB.Query("company").Select("*").Exec().Count()

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":        companiesBindDel,
		"total_count": totalCount,
		"page":        (offset / limit) + 1,
		"total_pages": totalPages,
	})
}

func FindCompany(c *gin.Context) {
	//Проверка наличия tip с указанным id в кэше
	id := c.Param("id")
	cacheKey := "company:" + id

	if cachedData, err := storage.Cache.Get(cacheKey); err == nil {
		var companyBind models_bind.CompanyBind
		if err := json.Unmarshal(cachedData, &companyBind); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": companyBind})
			return
		}
	}

	elem, found := storage.DB.Query("company").Where("id", reindexer.EQ, c.Param("id")).Get()

	if found {
		company := elem.(*models.Company)

		// Поиск подходящих man для данной company
		iterator := storage.DB.Query("man").Where("company_id", reindexer.EQ, company.Id).Sort("sort", true).Exec()
		defer iterator.Close()

		menBind := make([]models_bind.ManBind, 0)
		for iterator.Next() {
			man := iterator.Object().(*models.Man)

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
			menBind = append(menBind, manBind)
		}

		if err := iterator.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		companyBind := models_bind.CompanyBind{
			Id:          company.Id,
			Name:        company.Name,
			Established: company.Established,
			Men:         menBind,
		}

		// Кэширование данных
		jsonData, err := json.Marshal(companyBind)
		if err == nil {
			storage.Cache.Set(cacheKey, jsonData)
		}

		c.JSON(http.StatusOK, gin.H{"data": companyBind})
	} else {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error:": "Company with specified Id does not exist"})
		return
	}
}

func UpdateCompany(c *gin.Context) {
	elem, found := storage.DB.Query("company").
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

		if _, err := storage.DB.Update("company", item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": item})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Company with specified Id does not exist"})
		return
	}
}

func DeleteCompany(c *gin.Context) {
	elem, found := storage.DB.Query("company").
		Where("id", reindexer.EQ, c.Param("id")).
		Get()

	if found {
		company := elem.(*models.Company)

		// Удаление связанных man
		iteratorMan := storage.DB.Query("man").Where("company_id", reindexer.EQ, company.Id).Exec()
		defer iteratorMan.Close()

		for iteratorMan.Next() {
			man := iteratorMan.Object().(*models.Man)

			// Удаление связанных tip
			iteratorTip := storage.DB.Query("tip").Where("man_id", reindexer.EQ, man.Id).Exec()
			defer iteratorTip.Close()

			for iteratorTip.Next() {
				tip := iteratorTip.Object().(*models.Tip)
				if err := storage.DB.Delete("tip", tip); err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
			}

			if err := storage.DB.Delete("man", man); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		if err := storage.DB.Delete("company", company); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": company})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Company with specified Id does not exist"})
		return
	}
}
