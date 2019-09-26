package config

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	Port int // port to listen

	// Db-related config
	Dbhost     string
	Dbport     int
	Dbuser     string
	Dbpassword string
	Dbname     string
	Dbextra    string
}

func (c *Configuration) DbConnectionString() string {
	var result []string
	if c.Dbhost != "" {
		result = append(result, "host="+c.Dbhost)
	}
	if c.Dbport != 0 {
		result = append(result, "port="+strconv.Itoa(c.Dbport))
	}
	if c.Dbuser != "" {
		result = append(result, "user="+c.Dbuser)
	}
	if c.Dbpassword != "" {
		result = append(result, "password="+c.Dbpassword)
	}
	if c.Dbname != "" {
		result = append(result, "dbname="+c.Dbname)
	}
	if c.Dbextra != "" {
		result = append(result, strings.TrimSpace(c.Dbextra))
	}
	return strings.Join(result, " ")
}

func LoadConfig(path string) (Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return Configuration{}, err
	}
	var config Configuration
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}
