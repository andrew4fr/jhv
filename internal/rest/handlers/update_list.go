package handlers

import (
	"treasure/internal/rest/model"
	"treasure/internal/rest/server/operations"

	"github.com/go-openapi/runtime/middleware"
)

type ListUpdater interface {
	UpdateList() error
}

func UpdateList(u ListUpdater) operations.UpdateListHandlerFunc {
	return func(params operations.UpdateListParams) middleware.Responder {
		err := u.UpdateList()

		if err != nil {
			return operations.NewUpdateListServiceUnavailable().
				WithPayload(&model.Error{
					Code:   503,
					Result: "server error",
					Info:   err.Error(),
				})
		}

		return operations.NewUpdateListOK().
			WithPayload(&model.State{
				Code:   200,
				Result: "true",
			})
	}
}
