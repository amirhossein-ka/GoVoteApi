package mux

import (
	"net/http"
)

func (r *rest) routing() {
	r.router = r.router.StrictSlash(true)
	r.router.Use(r.handler.loggerMiddleware)
	api := r.router.PathPrefix("/api/v1").Subrouter()
	{
		user := api.PathPrefix("/user").Subrouter()
		user.HandleFunc("/register/", r.handler.register).Methods(http.MethodPost)
		user.HandleFunc("/delete/", r.handler.authorizationMiddleware(r.handler.delete)).Methods(http.MethodDelete)
		user.HandleFunc("/login/", r.handler.login).Methods(http.MethodPost)
		user.HandleFunc("/info/", r.handler.authorizationMiddleware(r.handler.info))
	}

	{
		// all requests to vote endpoints must be Authorized
		vote := api.PathPrefix("/vote").Subrouter()
		vote.Use(r.handler.authorizationMiddlewareMux)
	}
}
