package globals

import (
	"cant/util/config"
	"cant/util/livelog"
	"cant/util/stats"
	"database/sql"
)

// Globals
var (
	// version number
	Version = "0.0.5"
	// the users home path
	UserHomePath string
	// cant base pathe
	CantBasePath string
	// cant config path
	CantConfigPath string
	// the configuration
	Config *config.Config
	// database connection
	DB *sql.DB
	// global statistics reference
	Statistics = stats.New()
	// Livelog
	Livelog = livelog.New()
)
