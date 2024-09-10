package inits

import (
	"gostrecka/internal/service/database"
	"gostrecka/internal/service/database/sqlite"
	"log"

	"github.com/sarulabs/di/v2"
)

func InitDatabase(container di.Container) database.Database {
	var db database.Database
	var err error

	db = sqlite.New(container)
	err = db.Connect()

	if err != nil {
		log.Fatalf("Failed connecting to db %v", err)
	}

	return db
}
