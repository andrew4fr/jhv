package handlers

import (
	"treasure/internal/rest/model"
	"treasure/internal/rest/server/operations"

	"github.com/go-openapi/runtime/middleware"
)

type NamesGetter interface {
	GetNames(searchName, searchType string) (*model.Persons, error)
}

func GetNames(g NamesGetter) operations.GetNamesHandlerFunc {
	return func(params operations.GetNamesParams) middleware.Responder {
		names, err := g.GetNames(*params.Name, *params.Type)

		if err != nil {
			return operations.NewGetNamesServiceUnavailable().
				WithPayload(&model.Error{
					Code:   503,
					Result: "service unavailable",
					Info:   err.Error(),
				})
		}

		return operations.NewGetNamesOK().WithPayload(*names)
	}
}
