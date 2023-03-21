package handlers

import (
	"treasure/internal/rest/model"
	"treasure/internal/rest/server/operations"

	service "treasure"

	"github.com/go-openapi/runtime/middleware"
)

type StateGetter interface {
	GetState() *service.State
}

// ServiceState returns internal service state
func ServiceState(g StateGetter) operations.ServiceStateHandlerFunc {
	return func(params operations.ServiceStateParams) middleware.Responder {
		st := g.GetState()

		return operations.NewServiceStateOK().
			WithPayload(&model.State{
				Code:   int64(st.Code),
				Result: st.Result,
				Info:   st.Info,
			})
	}
}
