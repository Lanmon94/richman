package adaptor

import (
	"github.com/Lanmon94/richman/dbpool"
	"github.com/Lanmon94/richman/domain"
)

var (
	dbPool = dbpool.NewDBPool().Assert()
)

func init() {
	err := dbPool.NewConn().AutoMigrate(
		domain.FundStock{},
	).Error
	if nil != err {
		panic(err)
	}
}
