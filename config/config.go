package config

import (
	"errors"
	"fmt"
	"github.com/Lanmon94/richman/domain"
	"gopkg.in/yaml.v2"
	"os"
	"sync/atomic"
)

var (
	Config atomic.Value
)

func init() {
	local, err := loadLocationConf()
	if err != nil {
		panic(err)
	}
	var conf domain.DBConf
	conf = *local
	Config.Store(&conf)
}

func Conf() *domain.DBConf {
	return Config.Load().(*domain.DBConf)
}

//loadLocationConf load config from location config file
func loadLocationConf() (*domain.DBConf, error) {
	var conf domain.DBConf
	f, err := os.Open("conf.yml")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("loadLocationConfig path error :%s", err.Error()))
	}
	defer f.Close()
	if err = yaml.NewDecoder(f).Decode(&conf); err != nil {
		return nil, errors.New(fmt.Sprintf("loadLocationConfig decode error : %s", err.Error()))
	}
	return &conf, nil
}
