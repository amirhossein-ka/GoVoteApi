package mux

import (
	dto "GoVoteApi/DTO"
	"GoVoteApi/service/auth"
	voteSrv "GoVoteApi/service/vote"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *handler) createVote(w http.ResponseWriter, r *http.Request) {
	var req *dto.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(Error{Err: err, Section: "jsonDecoder", StatusCode: http.StatusBadRequest, Log: true, LogLevel: WarnLevel})
	}

	claim := r.Context().Value("claims")
	var claims = &auth.JwtClaims{}
	if claim != nil {
		claims = claim.(*auth.JwtClaims)
	}
	req.UserID = claims.ID

	voteResponse, err := h.srv.CreateVote(r.Context(), req)
	if err != nil {
		panic(Error{Err: err, Section: "createVote", StatusCode: http.StatusBadRequest, LogLevel: WarnLevel, Log: true})
	}
	//fmt.Println(voteResponse)
	writeJson(w, http.StatusCreated, voteResponse)
}

func (h *handler) getAllVotes(w http.ResponseWriter, r *http.Request) {
	// get query parameters
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		panic(Error{Err: err, Section: "jsonDecoder", StatusCode: http.StatusBadRequest, Log: true, LogLevel: WarnLevel})
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		panic(Error{Err: err, Section: "jsonDecoder", StatusCode: http.StatusBadRequest, Log: true, LogLevel: WarnLevel})
	}
	status, err := strconv.Atoi(r.URL.Query().Get("status"))
	if err != nil {
		panic(Error{Err: err, Section: "jsonDecoder", StatusCode: http.StatusBadRequest, Log: true, LogLevel: WarnLevel})
	}

	votes, err := h.srv.GetAllVotes(r.Context(), limit, offset, dto.VoteStatus(status))
	if err != nil {
		if err == voteSrv.ErrNoResult {
			panic(Error{Err: fmt.Errorf("no votes found"), Section: "getAllVotes/service get response",
				StatusCode: http.StatusNotFound, LogLevel: InfoLevel, Log: true})
		}
		panic(Error{Err: err, Section: "getAllVotes/repository", StatusCode: http.StatusInternalServerError, Log: true,
			LogLevel: WarnLevel})
	}

	writeJson(w, http.StatusFound, votes)
}

func (h *handler) getVoteWithID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	fmt.Println(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(Error{Err: err, Section: "getVoteWithID/convert id", StatusCode: http.StatusBadRequest, Log: true})
	}
	//fmt.Printf("moz: %d\n", uint(id))

	response, err := h.srv.GetVoteByID(r.Context(), uint(id))
	if err != nil {
		if err == sql.ErrNoRows {
			panic(Error{Err: fmt.Errorf("no vote found with given id"), Section: "getVoteWithID/getVote",
				StatusCode: http.StatusNotFound, LogLevel: InfoLevel, Log: true})
		}
		panic(Error{Err: err, Section: "getVoteWithID/getVote", StatusCode: http.StatusBadRequest, Log: true})
	}

	writeJson(w, http.StatusFound, response)
}

func (h *handler) getVoteWithSlug(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	response, err := h.srv.GetVoteBySlug(r.Context(), slug)
	if err != nil {
		if err == sql.ErrNoRows {
			panic(Error{Err: fmt.Errorf("no vote found with given slug"), Section: "getVoteWithSlug/service get response",
				StatusCode: http.StatusNotFound, LogLevel: InfoLevel, Log: true})
		}
		panic(Error{Err: err, Section: "getVoteWithSlug/service get response", StatusCode: http.StatusBadRequest,
			LogLevel: WarnLevel, Log: true})
	}

	writeJson(w, http.StatusFound, response)
}

func (h *handler) addVoteByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	voteID, err := strconv.Atoi(idStr)
	if err != nil {
		panic(Error{Err: err, Section: "getVoteWithID/convert id", StatusCode: http.StatusBadRequest, Log: true})
	}

	claim := r.Context().Value("claims")
	var claims = &auth.JwtClaims{}
	if claim != nil {
		claims = claim.(*auth.JwtClaims)
	}

	v := dto.Voters{}
	if err = json.NewDecoder(r.Body).Decode(&v); err != nil {
		panic(Error{Err: err, Section: "addVoteId/jsonDecoder", StatusCode: http.StatusBadRequest, LogLevel: WarnLevel, Log: true})
	}
	v.VoteID = uint(voteID)
	v.UserID = claims.ID
	v.Username = claims.Username

	response, err := h.srv.AddVote(r.Context(), &v)
	if err != nil {
		panic(Error{Err: err, Section: "addVoteID/service", StatusCode: http.StatusInternalServerError, LogLevel: WarnLevel, Log: true})
	}

	writeJson(w, http.StatusOK, response)
}

func (h *handler) addVoteBySlug(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	vote, err := h.srv.GetVoteBySlug(r.Context(), slug)
	if err != nil {
		panic(Error{Err: err, Section: "addVoteBySlug", StatusCode: http.StatusInternalServerError, LogLevel: WarnLevel, Log: true})
	}

	claim := r.Context().Value("claims")
	var claims = &auth.JwtClaims{}
	if claim != nil {
		claims = claim.(*auth.JwtClaims)
	}

	v := dto.Voters{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		panic(Error{Err: err, Section: "addVoteId/jsonDecoder", StatusCode: http.StatusBadRequest, LogLevel: WarnLevel, Log: true})
	}

	v.UserID = claims.ID
	v.VoteID = vote.ID
	v.Username = claims.Username
	response, err := h.srv.AddVoteSlug(r.Context(), slug, &v)
	if err != nil {
		panic(Error{Err: err, Section: "addVoteSlug/h.src.AddVote", StatusCode: http.StatusInternalServerError, LogLevel: WarnLevel, Log: true})
	}

	writeJson(w, http.StatusOK, response)
}

func (h *handler) addVoteOption(w http.ResponseWriter, r *http.Request) {

}
