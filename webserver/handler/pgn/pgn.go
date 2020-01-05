package pgn

import (
	pgnModel "cant/models/pgn"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllPGN returns all PGNs
func GetAllPGN(resp http.ResponseWriter, req *http.Request) {
	pgns, err := pgnModel.GetAll()
	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(resp)
	if err := enc.Encode(&pgns); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

// GetPGN returns a single PGN
func GetPGN(resp http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := pgnModel.GetByID(uint32(id))
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(resp)
	if err := enc.Encode(&res); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}
