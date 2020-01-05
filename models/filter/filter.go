package filter

import (
	"cant/util/globals"
	"database/sql"
)

// Filter represents a Filter without its implementation
type Filter struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	For  []ForPGN `json:"for_pgns"`
}

type ForPGN struct {
	PGN     uint32 `json:"pgn"`
	Enabled bool   `json:"enabled"`
}

// Filters is an array of Filter
type Filters []Filter

// New Filter
func New() Filter {
	return Filter{}
}

// prepared statements
var (
	idStmt        *sql.Stmt
	allStmt       *sql.Stmt
	insertStmt    *sql.Stmt
	insertForStmt *sql.Stmt
	pgnStmt       *sql.Stmt
	assocPGNStmt  *sql.Stmt
)

func Prepare() {
	if globals.DB != nil {
		allStmt, _ = globals.DB.Prepare("SELECT id FROM `filter`")
		idStmt, _ = globals.DB.Prepare("SELECT * FROM `filter` WHERE id=?")
		insertStmt, _ = globals.DB.Prepare(`INSERT INTO filter (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE id=VALUES(id), name=VALUES(name)`)
		insertForStmt, _ = globals.DB.Prepare(`INSERT INTO filter_for (filter, pgn, enabled) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE filter=VALUES(filter), pgn=VALUES(pgn), enabled=VALUES(enabled)`)
		pgnStmt, _ = globals.DB.Prepare("SELECT DISTINCT filter_for.filter FROM (cant.filter_for as filter_for, cant.filter as filt) where filter_for.filter=filt.id and (filter_for.pgn=? or filter_for.pgn=0)")
		assocPGNStmt, _ = globals.DB.Prepare("SELECT pgn, enabled FROM `filter_for` WHERE filter=?")
	}
}

// FilterCache offers a cached map for all filters
var FilterCache = map[int]Filter{}

// BuildCache builds a cache for all FilterCache in order to minimize stress on the GC
func BuildCache() error {
	filters, err := GetAll()
	if err != nil {
		return err
	}
	FilterCache = make(map[int]Filter)
	for i := range filters {
		FilterCache[filters[i].ID] = filters[i]
	}
	return nil
}

// GetAll fetches all Filter structs from the MySQL database by
func GetAll() (Filters, error) {
	rows, err := allStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	filters := Filters{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return filters, err
		}
		filter, err := GetByID(id)
		if err != nil {
			return filters, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

// GetByID fetches a Filter struct from the MySQL database by id
func GetByID(id int) (Filter, error) {
	row := idStmt.QueryRow(id)
	filter := New()
	if err := row.Scan(&filter.ID, &filter.Name); err != nil {
		return filter, err
	}
	rows, err := assocPGNStmt.Query(id)
	if err != nil {
		return filter, err
	}
	defer rows.Close()
	for rows.Next() {
		forPGN := ForPGN{}
		err := rows.Scan(&forPGN.PGN, &forPGN.Enabled)
		if err != nil {
			return filter, err
		}
		filter.For = append(filter.For, forPGN)
	}
	return filter, nil
}

// GetByPGN fetches a Filter struct from the MySQL database by id
func GetByPGN(id uint32) (Filters, error) {
	rows, err := pgnStmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	filters := Filters{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return filters, err
		}
		filter, err := GetByID(id)
		if err != nil {
			return filters, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

//Save or Update a given Filter in mysql
func (filter Filter) Save() error {
	if _, err := insertStmt.Exec(filter.ID, filter.Name); err != nil {
		return err
	}
	for _, pgn := range filter.For {
		if _, err := insertForStmt.Exec(filter.ID, pgn.PGN, pgn.Enabled); err != nil {
			return err
		}
	}
	return nil
}

// SaveAll saves all Filters to mysql
func (filters Filters) SaveAll() error {
	// start transaction
	tx, err := globals.DB.Begin()
	if err != nil {
		return err
	}
	stmt1, _ := tx.Prepare(`INSERT INTO filter (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE id=VALUES(id), name=VALUES(name)`)
	stmt2, _ := tx.Prepare(`INSERT INTO filter_for (filter, pgn, enabled) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE filter=VALUES(filter), pgn=VALUES(pgn), enabled=VALUES(enabled)`)
	defer stmt1.Close()
	defer stmt2.Close()
	for _, filter := range filters {
		if _, err := stmt1.Exec(filter.ID, filter.Name); err != nil {
			tx.Rollback()
			return err
		}
		for _, pgn := range filter.For {
			if _, err := stmt2.Exec(filter.ID, pgn.PGN, pgn.Enabled); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}
