package filter

import (
	filterModel "cant/models/filter"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllFilter returns all Filters
func GetAllFilter(resp http.ResponseWriter, req *http.Request) {
	filters, err := filterModel.GetAll()
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(resp)
	if err := enc.Encode(&filters); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

// GetFilterForPGN returns all available filters for a given pgn
func GetFilterForPGN(resp http.ResponseWriter, req *http.Request) {
	pgn, err := strconv.ParseUint(mux.Vars(req)["pgn"], 10, 32)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	res, err := filterModel.GetByPGN(uint32(pgn))
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

// DisableFilter disables the filter
func DisableFilter(resp http.ResponseWriter, req *http.Request) {
	pgn, err := strconv.ParseUint(mux.Vars(req)["pgn"], 10, 32)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	fid, err := strconv.Atoi(mux.Vars(req)["fid"])
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if _, ok := filterModel.FilterCache[fid]; !ok {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	for i := range filterModel.FilterCache[fid].For {
		if filterModel.FilterCache[fid].For[i].PGN == uint32(pgn) {
			filterModel.FilterCache[fid].For[i].Enabled = false
			filterModel.FilterCache[fid].Save()
			return
		}
	}

	resp.WriteHeader(http.StatusBadRequest)
}

// EnableFilter enables the filter
func EnableFilter(resp http.ResponseWriter, req *http.Request) {
	pgn, err := strconv.ParseUint(mux.Vars(req)["pgn"], 10, 32)
	if err != nil || pgn == 0 {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	fid, err := strconv.Atoi(mux.Vars(req)["fid"])
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if _, ok := filterModel.FilterCache[fid]; !ok {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	found := false
	zeroPresent := false
	for i := range filterModel.FilterCache[fid].For {
		if filterModel.FilterCache[fid].For[i].PGN == uint32(0) {
			zeroPresent = true
		}
		if filterModel.FilterCache[fid].For[i].PGN == uint32(pgn) {
			found = true
			filterModel.FilterCache[fid].For[i].Enabled = true
			filterModel.FilterCache[fid].Save()
			return
		}
	}
	if !found && zeroPresent {
		filt := filterModel.FilterCache[fid]
		filt.For = append(filt.For, filterModel.ForPGN{PGN: uint32(pgn), Enabled: true})
		filterModel.FilterCache[fid] = filt
		filterModel.FilterCache[fid].Save()
		return
	}

	resp.WriteHeader(http.StatusBadRequest)
}
