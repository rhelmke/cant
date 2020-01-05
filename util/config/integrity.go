package config

// this file adds integrity checks for a config.Config struct

import (
	"cant/util/database"
	"fmt"
	"net"
	"time"
)

// IntegrityCheck describes various functions used for checking the integrity of a config.Config struct
type IntegrityCheck struct {
	Name string
	Run  func() error
}

// AllIntegrityChecks returns an array containing all integrity checks
func (c *Config) AllIntegrityChecks() *[]*[]IntegrityCheck {
	return &[]*[]IntegrityCheck{
		c.MySQLIntegrityChecks(),
		c.WebinterfaceIntegrityChecks(),
	}
}

// MySQLIntegrityChecks returns a list of all mysql checks
func (c *Config) MySQLIntegrityChecks() *[]IntegrityCheck {
	// integrity checks are implemented as anonymous functions
	return &[]IntegrityCheck{
		IntegrityCheck{
			Name: "MySQL Connectivity",
			Run: func() error {
				db, err := database.NewMySQLConnection(c.MySQL.Host, c.MySQL.Port, c.MySQL.User, c.MySQL.Password, c.MySQL.DB)
				if err != nil {
					return err
				}
				defer db.Close()
				return nil
			},
		},
		IntegrityCheck{
			Name: "MySQL Database Existence",
			Run: func() error {
				db, err := database.NewMySQLConnection(c.MySQL.Host, c.MySQL.Port, c.MySQL.User, c.MySQL.Password, c.MySQL.DB)
				if err != nil {
					return err
				}
				defer db.Close()
				var schemaName string
				if err = db.QueryRow("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME=?", c.MySQL.DB).Scan(&schemaName); err != nil {
					return err
				}
				return nil
			},
		},
		IntegrityCheck{
			Name: "MySQL Permissions",
			Run: func() error {
				db, err := database.NewMySQLConnection(c.MySQL.Host, c.MySQL.Port, c.MySQL.User, c.MySQL.Password, c.MySQL.DB)
				if err != nil {
					return err
				}
				defer db.Close()
				rndTableName := "cant-" + time.Now().String()
				stmts := []string{
					"CREATE TABLE `" + rndTableName + "`(`id` INT NOT NULL, `text` TEXT, PRIMARY KEY (ID))",
					"INSERT INTO `" + rndTableName + "` (`id`, `text`) VALUES ('1', 'a')",
					"SELECT * FROM `" + rndTableName + "`",
					"UPDATE `" + rndTableName + "` SET `text`='b'",
					"DELETE FROM `" + rndTableName + "` WHERE `id`='1'",
					"DROP TABLE `" + rndTableName + "`",
				}
				for _, stmt := range stmts {
					_, err = db.Exec(stmt)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
	}
}

// WebinterfaceIntegrityChecks returns a list of all webinterface checks
func (c *Config) WebinterfaceIntegrityChecks() *[]IntegrityCheck {
	return &[]IntegrityCheck{
		IntegrityCheck{
			Name: "Webinterface Bind address",
			Run: func() error {
				listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.Webinterface.Host, c.Webinterface.Port))
				if err != nil {
					return err
				}
				defer listener.Close()
				return nil
			},
		},
	}
}
