package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sync"
)

var cfg *Config
var once sync.Once

type Config struct {
	MongoConf map[string]MongoConfig `yaml:"mongo_conf"`
}

type MongoConfig struct {
	Addr            string `yaml:"addr"`
	PoolMaxSize     int    `yaml:"pool_max_size"`
	PoolMinSize     int    `yaml:"pool_min_size"`
	MaxConnIdleTime int    `yaml:"max_conn_idle_time"`
}

func NewConfig() *Config {
	return &Config{}
}

func GetInstance() *Config {
	once.Do(func() {
		cfg = NewConfig()
	})
	return cfg
}

func (cfg *Config) InitConfig(file string) error {
	conf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("read config file failed err = %s\n", err)
		return err
	}

	if err := yaml.Unmarshal(conf, cfg); err != nil {
		fmt.Printf("parse config file failed err = %s\n", err)
		return err
	}
	return err
}
