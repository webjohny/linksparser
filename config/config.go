package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Env      string
	Port      string
	ApiHost   string
	ApiLink   string
	ApiKey    string
	ApiSecret string
	MysqlHost  string
	MysqlDb   string
	MysqlLogin   string
	MysqlPass   string
	ImgPath   string
}

func (conf *Configuration) Create(filename string) {

	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&conf)
}