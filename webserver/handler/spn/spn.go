package spn

import (
	spnModel "cant/models/spn"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllSPN returns all SPNs
func GetAllSPN(resp http.ResponseWriter, req *http.Request) {
	spns, err := spnModel.GetAll()
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(resp)
	if err := enc.Encode(&spns); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

// GetSPN returns a single SPN
func GetSPN(resp http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	res, err := spnModel.GetByID(uint32(id))
	if err != nil {
		fmt.Println(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(resp)
	if err := enc.Encode(&res); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

func GetSPNsForPGN(resp http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 32)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := spnModel.GetByPGN(uint32(id))
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(resp)
	if err := enc.Encode(&res); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}
