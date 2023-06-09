// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ServiceStateHandlerFunc turns a function with the right signature into a service state handler
type ServiceStateHandlerFunc func(ServiceStateParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ServiceStateHandlerFunc) Handle(params ServiceStateParams) middleware.Responder {
	return fn(params)
}

// ServiceStateHandler interface for that can handle valid service state params
type ServiceStateHandler interface {
	Handle(ServiceStateParams) middleware.Responder
}

// NewServiceState creates a new http.Handler for the service state operation
func NewServiceState(ctx *middleware.Context, handler ServiceStateHandler) *ServiceState {
	return &ServiceState{Context: ctx, Handler: handler}
}

/* ServiceState swagger:route GET /state serviceState

ServiceState service state API

*/
type ServiceState struct {
	Context *middleware.Context
	Handler ServiceStateHandler
}

func (o *ServiceState) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewServiceStateParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
