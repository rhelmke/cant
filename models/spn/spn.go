package spn

import (
	"cant/util/globals"
	"database/sql"
)

// SPN represents a Suspect Parameter Number
type SPN struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	PGN  uint32 `json:"pgn"`
}

// SPNs is an array of *SPN
type SPNs []SPN

// New SPN
func New() SPN {
	return SPN{}
}

// prepared statements
var (
	allStmt    *sql.Stmt
	idStmt     *sql.Stmt
	insertStmt *sql.Stmt
	pgnStmt    *sql.Stmt
)

func Prepare() {
	if globals.DB != nil {
		allStmt, _ = globals.DB.Prepare("SELECT * FROM `spn`")
		idStmt, _ = globals.DB.Prepare("SELECT * FROM `spn` WHERE id=?")
		insertStmt, _ = globals.DB.Prepare(`INSERT INTO spn (id, name, pgn) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id=VALUES(id), name=VALUES(name), pgn=VALUES(pgn)`)
		pgnStmt, _ = globals.DB.Prepare("SELECT * FROM `spn` WHERE pgn=?")
	}
}

// GetByID fetches a SPN struct from the MySQL database by id
func GetByID(id uint32) (SPN, error) {
	row := idStmt.QueryRow(id)
	spn := New()
	if err := row.Scan(&spn.ID, &spn.Name, &spn.PGN); err != nil {
		return spn, err
	}
	return spn, nil
}

// GetByPGN fetches a SPN struct from the MySQL database by id
func GetByPGN(id uint32) (SPNs, error) {
	rows, err := pgnStmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	spns := SPNs{}
	for rows.Next() {
		spn := New()
		err := rows.Scan(&spn.ID, &spn.Name, &spn.PGN)
		if err != nil {
			return spns, err
		}
		spns = append(spns, spn)
	}
	return spns, nil
}

// GetAll fetches all SPN structs from the MySQL database
func GetAll() (SPNs, error) {
	rows, err := allStmt.Query()
	if err != nil {
		return nil, err
	}
	spns := SPNs{}
	for rows.Next() {
		spn := New()
		err := rows.Scan(&spn.ID, &spn.Name, &spn.PGN)
		if err != nil {
			return spns, err
		}
		spns = append(spns, spn)
	}
	rows.Close()
	return spns, nil
}

//Save or Update a given SPN in mysql
func (spn SPN) Save() error {
	if _, err := insertStmt.Exec(spn.ID, spn.Name, spn.PGN); err != nil {
		return err
	}
	return nil
}

// SaveAll saves all SPNs to mysql
func (spns SPNs) SaveAll() error {
	// start transaction
	tx, err := globals.DB.Begin()
	if err != nil {
		return err
	}
	stmt, _ := tx.Prepare(`INSERT INTO spn (id, name, pgn) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id=VALUES(id), name=VALUES(name), pgn=VALUES(pgn)`)
	defer stmt.Close()
	for _, spn := range spns {
		if _, err := stmt.Exec(spn.ID, spn.Name, spn.PGN); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
