package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto" // use Reindexer as standalone server and connect to it via TCP.

	"github.com/Maritornez/GoCRUD/internal/config"
	"github.com/Maritornez/GoCRUD/internal/models"
)

var DB *reindexer.Reindexer
var Cache *bigcache.BigCache

func ConnectDatabase() {
	config, err_config := config.NewConfig()
	if err_config != nil {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Println("Current working directory:", dir)
		panic(err_config)
	}

	var dsn = "cproto://" + config.Database.User + ":" + config.Database.Pass +
		"@" + config.Database.IpAddress + ":" + config.Database.Port + "/" + config.Database.DBName

	// Init a database instance and choose the binding (connect to server)
	// Database should be created explicitly via reindexer_tool or via WithCreateDBIfMissing option:
	// If server security mode is enabled, then username and password are mandatory
	DB = reindexer.NewReindex(dsn, reindexer.WithCreateDBIfMissing())

	// Проверка, была ли инициализирована БД. То есть проверка подключения к БД
	if DB.Status().Err == nil {
		fmt.Println("Connected to DB")

		// Открывите пространств имен. Если пространство не существует, то оно будет создано
		if err := DB.OpenNamespace("man", reindexer.DefaultNamespaceOptions(), models.Man{}); err != nil {
			panic(err)
		}
		if err := DB.OpenNamespace("tip", reindexer.DefaultNamespaceOptions(), models.Tip{}); err != nil {
			panic(err)
		}
		if err := DB.OpenNamespace("company", reindexer.DefaultNamespaceOptions(), models.Company{}); err != nil {
			panic(err)
		}

	}

	//DB.DropNamespace("man") // Comment!!!!!
}

func CloseDatabase() {
	DB.Close()
}

func InitCache() {
	Cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(15*time.Minute))
}

func InitializeDatabase() {
	// Проверка, есть ли в базе данных данные
	iteratorCom := DB.Query("company").Select("*").Limit(1).Exec()
	defer iteratorCom.Close()
	iteratorMan := DB.Query("man").Select("*").Limit(1).Exec()
	defer iteratorMan.Close()
	iteratorTip := DB.Query("tip").Select("*").Limit(1).Exec()
	defer iteratorTip.Close()

	if !(iteratorTip.Next() || iteratorMan.Next() || iteratorCom.Next()) {

		initialDataCompanies := []models.Company{
			{Id: 1, Name: "Pyaterochka", Established: 2010},
			{Id: 2, Name: "Gazprom", Established: 1989},
			{Id: 3, Name: "Sberbank", Established: 1841},
		}

		for _, company := range initialDataCompanies {
			if err := DB.Upsert("company", company); err != nil {
				fmt.Println("Ошибка при добавлении данных company:", err)
			}
		}

		initialDataMans := []models.Man{
			{Id: 1, Name: "Alex", Age: 30, CompanyId: 1, Sort: 1},
			{Id: 2, Name: "Ivan", Age: 24, CompanyId: 1, Sort: 3},
			{Id: 3, Name: "Max", Age: 27, CompanyId: 1, Sort: 5},
			{Id: 4, Name: "Petr", Age: 29, CompanyId: 2, Sort: 2},
			{Id: 5, Name: "Vlad", Age: 22, CompanyId: 2, Sort: 4},
		}

		for _, man := range initialDataMans {
			if err := DB.Upsert("man", man); err != nil {
				fmt.Println("Ошибка при добавлении данных man:", err)
			}
		}

		initialDataTips := []models.Tip{
			{
				Id: 1, ManId: 1, Title: "Credits",
				Pages: []models.Page{
					{Title: "Sberbank", Content: "100000"},
					{Title: "Alfa Bank", Content: "150000"},
				},
			},
			{
				Id: 2, ManId: 1, Title: "Hobby",
				Pages: []models.Page{
					{Title: "Piano", Content: "From 10 years old"},
				},
			},
			{
				Id: 3, ManId: 3, Title: "Weight",
				Pages: []models.Page{
					{Title: "Weight", Content: "75 kg"},
				},
			},
			{
				Id: 4, ManId: 4, Title: "Accounts",
				Pages: []models.Page{
					{Title: "Steam", Content: "Petr1995"},
					{Title: "Google", Content: "pertIvanov@gmail.com"},
				},
			},
		}

		for _, tip := range initialDataTips {
			if err := DB.Upsert("tip", tip); err != nil {
				fmt.Println("Ошибка при добавлении данных tip:", err)
			}
		}
	}
}
