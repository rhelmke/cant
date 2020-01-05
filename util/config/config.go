// Package config implements a basic structure for json-marshalled configuration files
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Config is the json root, containing multiple subobjects
type Config struct {
	MySQL        dbConfig  `json:"mysql"`
	Webinterface webConfig `json:"webinterface"`
	Network      netConfig `json:"network"`
}

// dbConfigData represents the mysql connection information
type dbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

// webConfig represents the webserver configuration
type webConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// netConfig represents the network configuration
type netConfig struct {
	Interface0 string `json:"interface0"`
	Interface1 string `json:"interface1"`
}

// NewConfig creates a new config
func NewConfig() *Config {
	return &Config{}
}

// ToJSON can be called on a config object and marshals the struct to binary json
func (c *Config) ToJSON() ([]byte, error) {
	res, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return []byte{}, err
	}
	return res, nil
}

// Save the configuration to a json file
func (c *Config) Save(file string) error {
	blob, err := c.ToJSON()
	if err != nil {
		return err
	}
	_, err = os.Stat(path.Dir(file))
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path.Dir(file), 0700); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return ioutil.WriteFile(file, blob, 0600)
}
