package models

import (
	"fmt"

	"github.com/restream/reindexer/v3"
	// use Reindexer as standalone server and connect to it via TCP.
	_ "github.com/restream/reindexer/v3/bindings/cproto"

	"github.com/Maritornez/GoCRUD/internal/config"
)

var DB *reindexer.Reindexer

func ConnectDatabase() {
	config, err_config := config.NewConfig(config.YamlPath)
	if err_config != nil {
		panic(err_config)
	}

	var dsn = "cproto://" + config.Database.User + ":" + config.Database.Pass +
		"@" + config.Database.Host + ":" + config.Database.Port + "/" + config.Database.DBName

	// Init a database instance and choose the binding (connect to server)
	// Database should be created explicitly via reindexer_tool or via WithCreateDBIfMissing option:
	// If server security mode is enabled, then username and password are mandatory
	DB = reindexer.NewReindex(dsn, reindexer.WithCreateDBIfMissing())

	// Check if DB was initialized correctly
	if DB.Status().Err != nil {
		panic(DB.Status().Err)
	} else {
		fmt.Println("Connected to DB")
	}

	//DB.DropNamespace("mans") // Comment!!!!!

	// Create new namespace (if there's not) with name 'mans', which will store structs of type 'Item'
	if err := DB.OpenNamespace("mans", reindexer.DefaultNamespaceOptions(), Man{}); err != nil {
		panic(err)
	}
}

func CloseDatabase() {
	DB.Close()
}
