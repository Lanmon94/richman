package dbpool

import (
	"fmt"
	"github.com/Lanmon94/richman/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Client struct {
	DB *gorm.DB
}

// NewClient Create a MySQL client, retry if failed
func NewClient(conf *domain.DBConf) *Client {
	var err error
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("mysql Connect Error:", err)
		db = nil
		return nil
	}
	return &Client{db}
}

// NewConn Fetch a db connection from pool
func (c *Client) NewConn() *gorm.DB {
	db := c.DB.New()
	db.BlockGlobalUpdate(true)
	return db
}
