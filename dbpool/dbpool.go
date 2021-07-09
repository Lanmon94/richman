package dbpool

import (
	"errors"
	"fmt"
	"github.com/Lanmon94/richman/config"
	"github.com/jinzhu/gorm"
	"sync"
)

var (
	lock sync.Mutex
)

type DBPool struct {
	pool *Client
	lock sync.Mutex
}

func InitDB() (*Client, error) {
	lock.Lock()
	defer lock.Unlock()
	dbConf := config.Conf()
	pool := NewClient(dbConf)
	return pool, nil
}

func (p *DBPool) NewConn() *gorm.DB {
	if p.pool != nil {
		return p.pool.NewConn()
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.pool != nil {
		return p.pool.NewConn()
	}
	var err error
	p.pool, err = InitDB()
	if err != nil {
		panic(errors.New(fmt.Sprintf("Init db[%s]", err.Error())))
	}
	return p.pool.NewConn()
}

func (p *DBPool) Assert() *DBPool {
	p.NewConn()
	return p
}

func NewDBPool() *DBPool {
	return &DBPool{}
}
