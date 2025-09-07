package initialization

import (
	"arch/internal/application"
	"arch/internal/delivery/http"
	"arch/internal/delivery/middleware"
	"arch/internal/infrastructure"

	"github.com/jmoiron/sqlx"
)

func InitLayers() (delivery *http.Http, businessDatabase *sqlx.DB) {

	businessDatabase = infrastructure.NewBusinessDatabase(ConfigService)
	sources := infrastructure.Sources{
		BusinessDB: businessDatabase,
	}
	infrastructures := infrastructure.New(&sources)
	app := application.New(infrastructures)
	middleware := middleware.New()
	delivery = http.NewHttp(app, middleware)
	return delivery, businessDatabase
}
