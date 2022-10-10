package mux

import (
	"GoVoteApi/config"
	"GoVoteApi/controller"
	"GoVoteApi/pkg/logger"
	"GoVoteApi/service"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type handler struct {
	srv service.Service
	cfg config.Secrets
}

type rest struct {
	handler *handler
	server  *http.Server
	router  *mux.Router
	cfg     config.Server
	recover *recovery
}

func (r *rest) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer r.recover.recoverHttp(w)
	println(req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	r.router.ServeHTTP(w, req)
}

// Start implements controller.Rest
func (r *rest) Start(addr string) error {
	log.Println("Starting server on:", addr)
	r.routing()
	r.server = &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	return r.server.ListenAndServe()
}

// Stop implements controller.Rest
func (r *rest) Stop() error {
	log.Println("Stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return r.server.Shutdown(ctx)
}

func New(srv service.Service, cfg *config.Config, logger logger.Logger) controller.Rest {
	return &rest{
		router:  mux.NewRouter(),
		cfg:     cfg.Server,
		recover: &recovery{logger: logger},
		handler: &handler{
			srv: srv,
			cfg: cfg.Secrets,
		},
	}
}
