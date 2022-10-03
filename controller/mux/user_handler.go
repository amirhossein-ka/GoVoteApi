package mux

import (
	dto "GoVoteApi/DTO"
	"GoVoteApi/service/auth"
	"GoVoteApi/service/user"
	"encoding/json"
	"net/http"
)

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, dto.Error{Status: dto.StatusError, Error: err.Error()})
		return
	}

	resp, err := h.srv.Register(r.Context(), &req)
	if err != nil {
		if err == user.ErrUsernameExists {
			writeJson(w, http.StatusBadRequest, dto.Error{
				Status: dto.StatusError,
				Error:  err.Error(),
			})
			return
		}
		writeJson(w, http.StatusInternalServerError, dto.Error{
			Status: dto.StatusError,
			Error:  err.Error(),
		})
		return
	}

	writeJson(w, http.StatusCreated, resp)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, dto.Error{Status: dto.StatusError, Error: err.Error()})
		return
	}

	resp, err := h.srv.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeJson(w, http.StatusBadRequest, dto.Error{Status: dto.StatusError, Error: err.Error()})
		return
	}
	writeJson(w, http.StatusCreated, resp)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	claim := r.Context().Value("claims")
	var claims = &auth.JwtClaims{}
	if claim != nil {
		claims = claim.(*auth.JwtClaims)
	}

	if err := h.srv.Delete(r.Context(), &dto.UserRequest{ID: claims.ID}); err != nil {
		writeJson(w, http.StatusBadRequest, dto.Error{Status: dto.StatusError, Error: err.Error()})
		return
	}

	writeJson(w, http.StatusOK, &dto.UserResponse{Status: dto.StatusDeleted, Data: "user deleted"})
}

func (h *handler) info(w http.ResponseWriter, r *http.Request) {
	claim := r.Context().Value("claims")
	var claims = &auth.JwtClaims{}
	if claim != nil {
		claims = claim.(*auth.JwtClaims)
	}

	username := claims.Username

	u, err := h.srv.Info(r.Context(), username)
	if err != nil {
		writeJson(w, http.StatusBadRequest, dto.Error{Status: dto.StatusError, Error: err.Error()})
		return
	}

	writeJson(w, http.StatusFound, u)
}
