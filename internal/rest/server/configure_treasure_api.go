// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	_ "embed"
	"net/http"

	"database/sql"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	service "treasure"
	"treasure/internal/rest/handlers"
	"treasure/internal/rest/server/operations"
	"treasure/storage"
)

//go:generate swagger generate server --target ../../../../jhv --name TreasureAPI --spec ../../../api/swagger.yaml --model-package internal/rest/model --server-package internal/rest/server --principal interface{}

//go:embed postgres_schema.sql
var dbSchema string

func configureFlags(api *operations.TreasureAPIAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Service configuration",
			Options:          &globalConfig,
		},
	}
}

func configureAPI(api *operations.TreasureAPIAPI) http.Handler {
	api.ServeError = errors.ServeError
	api.UseSwaggerUI()
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	//db, err := sql.Open("mysql", globalConfig.StorageDSN)
	db, err := sql.Open("pgx", globalConfig.StorageDSN)
	if err != nil {
		panic("error on create storage connection")
	}
	log.Debug().Msg("successful db conn")

	_, err = db.Exec(dbSchema)
	if err != nil {
		panic("error on create db schema")
	}

	stor := storage.NewPostgresStorage(db)
	srv := service.New(globalConfig.XMLPath, stor)

	api.GetNamesHandler = handlers.GetNames(srv)
	api.ServiceStateHandler = handlers.ServiceState(srv)
	api.UpdateListHandler = handlers.UpdateList(srv)

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func configureTLS(tlsConfig *tls.Config) {
}

func configureServer(s *http.Server, scheme, addr string) {
}

func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
