package config

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		Name     string `yaml:"name"`
		Env      string `yaml:"env"`
		Debug    bool   `yaml:"debug"`
		Timezone string `yaml:"timezone"`
	} `yaml:"app"`
	Redis struct {
		Host                     string `yaml:"host"`
		Db                       int    `yaml:"db"`
		FrontierInquiryRetention int    `yaml:"frontier_inquiry_retention"`
		CachePrefix              string `yaml:"cache_prefix"`
		EncryptionKey            string `yaml:"encryption_key"`
	} `yaml:"redis"`
	Mongo struct {
		DatabaseName string   `yaml:"database_name"`
		Host         []string `yaml:"host"`
		Username     string   `yaml:"username"`
		Password     string   `yaml:"password"`
	} `yaml:"mongo"`
	Serve struct {
		Port         int `yaml:"port"`
		WriteTimeout int `yaml:"write_timeout"`
		ReadTimeout  int `yaml:"read_timeout"`
	} `yaml:"serve"`
	FCM struct {
		ApiKey string `yaml:"api_key"`
	} `yaml:"fcm"`
}

func Load(pathConfig string) Config {
	var config Config
	filename, _ := filepath.Abs(pathConfig)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic("load config file", err.Error())
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Panic("parsing config file", err.Error())
	}
	return config
}
