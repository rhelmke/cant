// Package pgn implements the database PGN model
package pgn

import (
	"cant/util/globals"
	"database/sql"
)

// PGN represents a Parameter Group Number
type PGN struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	EDP         int    `json:"edp"`
	DP          int    `json:"dp"`
	PF          int    `json:"pf"`
	PS          string `json:"ps"`
	Multipacket bool   `json:"multipacket"`
	DLC         int    `json:"dlc"`
}

// prepared statements
var (
	allStmt    *sql.Stmt
	idStmt     *sql.Stmt
	insertStmt *sql.Stmt
)

// PGNs is an array of PGN
type PGNs []PGN

// New PGN
func New() PGN {
	return PGN{}
}

// Prepare the sql statements
func Prepare() {
	if globals.DB != nil {
		allStmt, _ = globals.DB.Prepare("SELECT * FROM `pgn`")
		idStmt, _ = globals.DB.Prepare("SELECT * FROM `pgn` WHERE id=?")
		insertStmt, _ = globals.DB.Prepare(`INSERT INTO pgn (id, name, edp, dp, pf, ps, multipacket, pgn_dlc) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=VALUES(id), name=VALUES(name), edp=VALUES(edp), dp=VALUES(dp), pf=VALUES(pf), ps=VALUES(ps), multipacket=VALUES(multipacket), pgn_dlc=VALUES(pgn_dlc)`)
	}
}

// PGNCache offers a cached map for all pgns
var PGNCache = map[uint32]PGN{}

// BuildCache builds a cache for all PGNs in order to minimize stress on the GC
func BuildCache() error {
	pgns, err := GetAll()
	if err != nil {
		return err
	}
	PGNCache = make(map[uint32]PGN)
	for i := range pgns {
		PGNCache[pgns[i].ID] = pgns[i]
	}
	return nil
}

// GetByID fetches a PGN struct from the MySQL database by id
func GetByID(id uint32) (PGN, error) {
	if _, ok := PGNCache[id]; ok {
		return PGNCache[id], nil
	}
	pgn := New()
	row := idStmt.QueryRow(id)
	if err := row.Scan(&pgn.ID, &pgn.Name, &pgn.EDP, &pgn.DP, &pgn.PF, &pgn.PS, &pgn.Multipacket, &pgn.DLC); err != nil {
		return pgn, err
	}
	return pgn, nil
}

// GetAll fetches all PGN structs from the MySQL database
func GetAll() (PGNs, error) {
	rows, err := allStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pgns := PGNs{}
	for rows.Next() {
		pgn := New()
		if err := rows.Scan(&pgn.ID, &pgn.Name, &pgn.EDP, &pgn.DP, &pgn.PF, &pgn.PS, &pgn.Multipacket, &pgn.DLC); err != nil {
			return nil, err
		}
		pgns = append(pgns, pgn)
	}
	return pgns, nil
}

// Save or Update a given PGN in mysql
func (pgn PGN) Save() error {
	if _, err := insertStmt.Exec(pgn.ID, pgn.Name, pgn.EDP, pgn.DP, pgn.PF, pgn.PS, pgn.Multipacket, pgn.DLC); err != nil {
		return err
	}
	return nil
}

// SaveAll saves all PGNs to mysql
func (pgns PGNs) SaveAll() error {
	// start transaction
	tx, err := globals.DB.Begin()
	if err != nil {
		return err
	}
	stmt, _ := tx.Prepare(`INSERT INTO pgn (id, name, edp, dp, pf, ps, multipacket, pgn_dlc) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=VALUES(id), name=VALUES(name), edp=VALUES(edp), dp=VALUES(dp), pf=VALUES(pf), ps=VALUES(ps), multipacket=VALUES(multipacket), pgn_dlc=VALUES(pgn_dlc)`)
	defer stmt.Close()
	for _, pgn := range pgns {
		if _, err := stmt.Exec(pgn.ID, pgn.Name, pgn.EDP, pgn.DP, pgn.PF, pgn.PS, pgn.Multipacket, pgn.DLC); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
