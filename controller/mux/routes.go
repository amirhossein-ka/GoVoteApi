package mux

import (
	"net/http"
)

//
// TODO: BUG => it returns null and 200 when you vote twice on type 1 (which is quiz)

func (r *rest) routing() {
	r.router = r.router.StrictSlash(true)
	r.router.Use(r.handler.loggerMiddleware)
	api := r.router.PathPrefix("/api/v1").Subrouter()
	{
		user := api.PathPrefix("/user").Subrouter()
		user.HandleFunc("/register/", r.handler.register).Methods(http.MethodPost)
		user.HandleFunc("/delete/", r.handler.authorizationMiddleware(r.handler.delete)).Methods(http.MethodDelete)
		user.HandleFunc("/login/", r.handler.login).Methods(http.MethodPost)
		user.HandleFunc("/info/", r.handler.authorizationMiddleware(r.handler.info)).Methods(http.MethodGet)
	}

	{
		vote := api.PathPrefix("/vote").Subrouter()
		// all requests to vote endpoints must be Authorized
		vote.Use(r.handler.authorizationMiddlewareMux)
		vote.HandleFunc("/create/", r.handler.createVote).Methods(http.MethodPost)
		vote.HandleFunc("/add/{id:[0-9]+}/", r.handler.addVoteByID).Methods(http.MethodPost)
		vote.HandleFunc("/add/{slug}/", r.handler.addVoteBySlug).Methods(http.MethodPost)

		vote.HandleFunc("/all/", r.handler.getAllVotes).Methods(http.MethodGet)
		vote.HandleFunc("/{id:[0-9]+}/", r.handler.getVoteWithID).Methods(http.MethodGet)
		vote.HandleFunc("/{slug}/", r.handler.getVoteWithSlug).Methods(http.MethodGet)
	}
}
