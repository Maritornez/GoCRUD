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
	config, err_config := config.NewConfig(config.YamlPath)
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
