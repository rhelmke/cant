// Package database provides a wrapper for database/sql and the corresponding mysql driver
package database

import (
	"database/sql"
	"fmt"

	// mysql driver will be injected using reflection. No reason to pollute the namespace
	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLConnection creates a new MySQL Connection
func NewMySQLConnection(host string, port int, user string, password string, database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		defer db.Close()
		return nil, err
	}
	return db, nil
}

// ClearDatabase removes all tables from a database
func ClearDatabase(db *sql.DB) error {
	rows, err := db.Query("SELECT GROUP_CONCAT('`', table_name, '`') FROM information_schema.tables WHERE table_schema = (SELECT DATABASE())")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil
		}
		for _, stmt := range []string{"SET FOREIGN_KEY_CHECKS = 0", "DROP TABLE IF EXISTS " + table, "SET FOREIGN_KEY_CHECKS = 1"} {
			if _, err := db.Exec(stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

// WipeTable wipes a given table
func WipeTable(db *sql.DB, table string) error {
	for _, stmt := range []string{"SET FOREIGN_KEY_CHECKS = 0", "TRUNCATE TABLE " + table, "SET FOREIGN_KEY_CHECKS = 1"} {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

// CreateTables creates all tables needed for cant
func CreateTables(db *sql.DB) error {
	if err := ClearDatabase(db); err != nil {
		fmt.Println(err)
		return err
	}
	stmts := []string{
		`CREATE TABLE pgn (id INT unsigned NOT NULL,
                           name text,
                           edp INT DEFAULT NULL,
                           dp INT DEFAULT NULL,
                           pf INT DEFAULT NULL,
                           ps TEXT DEFAULT NULL,
                           multipacket BOOLEAN DEFAULT NULL,
                           pgn_dlc INT DEFAULT NULL,
                           PRIMARY KEY (id),
                           UNIQUE KEY id_UNIQUE (id))`,
		`CREATE TABLE spn (id INT unsigned NOT NULL,
                           name TEXT,
                           pgn INT unsigned DEFAULT NULL,
                           PRIMARY KEY (id),
                           UNIQUE KEY id_UNIQUE (id),
                           KEY pgn_idx (pgn),
                           CONSTRAINT pgn
                                FOREIGN KEY (pgn)
                                REFERENCES pgn (id)
                                ON DELETE NO ACTION 
                                ON UPDATE NO ACTION)`,
		`CREATE TABLE filter (id int NOT NULL,
                              name text NOT NULL,
                              PRIMARY KEY (id))`,
		`CREATE TABLE filter_for (filter int NOT NULL,
                                  pgn int unsigned DEFAULT NULL,
                                  UNIQUE(filter, pgn),
                                  enabled boolean DEFAULT NULL,
                                  CONSTRAINT filter 
                                    FOREIGN KEY (filter)
                                    REFERENCES filter (id)
                                    ON DELETE NO ACTION 
                                    ON UPDATE NO ACTION,
                                  CONSTRAINT pgn_filter 
                                    FOREIGN KEY (pgn)
                                    REFERENCES pgn (id)
                                    ON DELETE NO ACTION 
                                    ON UPDATE NO ACTION)`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
