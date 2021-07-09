package adaptor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Lanmon94/bar"
	"github.com/Lanmon94/richman/domain"
	"github.com/Lanmon94/richman/util/http"
	"github.com/robertkrimen/otto"
	"strings"
)

func InitFundData() {
	allFund := GetAllFund()
	ub := bar.NewBar()
	ub.NewOption(0, int64(len(allFund)))
	for index, fund := range allFund {
		ub.Play(int64(index))
		stockCodesData, syl1y, syl3y, syl6y := GetFundSyl(fund[0])
		stockCodes := strings.Split(stockCodesData, ",")
		for _, stockcode := range stockCodes {
			c := strings.Split(stockcode, ".")
			code := "0"
			if len(c) == 2 {
				code = c[1]
			}
			if len(c) == 1 {
				code = c[0]
			}
			dbPool.NewConn().Model(&domain.FundStock{}).Create(&domain.FundStock{
				StockCode: code,
				FundCode:  fund[0],
				FundName:  fund[2],
				Syl1y:     syl1y,
				Syl3y:     syl3y,
				Syl6y:     syl6y,
			})
		}
	}
	ub.Finish()
}

type Fund [][]string

func GetAllFund() Fund {
	res, _ := http.Get("http://fund.eastmoney.com/js/fundcode_search.js")
	body := bytes.TrimPrefix(res, []byte("\xef\xbb\xbf"))
	first := strings.Replace(string(body), "var r = ", "", -1)
	end := strings.TrimSuffix(first, ";")
	var d Fund
	json.Unmarshal([]byte(end), &d)
	return d
}

func GetFundSyl(code string) (string, float64, float64, float64) {
	var (
		stockCodes string  = ""
		syl1y      float64 = 0
		syl3y      float64 = 0
		syl6y      float64 = 0
	)
	res, err := http.Get(fmt.Sprintf("http://fund.eastmoney.com/pingzhongdata/"+"%s", code+".js"))
	if err != nil {
		return stockCodes, syl1y, syl3y, syl6y
	}
	vm := otto.New()
	vm.Run(string(res))
	if value, err := vm.Get("stockCodesNew"); err == nil {
		if value_int, err := value.ToString(); err == nil {
			stockCodes = value_int
		}
	}
	if value, err := vm.Get("syl_1y"); err == nil {
		if value_int, err := value.ToFloat(); err == nil {
			syl1y = value_int
		}
	}
	if value, err := vm.Get("syl_3y"); err == nil {
		if value_int, err := value.ToFloat(); err == nil {
			syl3y = value_int
		}
	}
	if value, err := vm.Get("syl_6y"); err == nil {
		if value_int, err := value.ToFloat(); err == nil {
			syl6y = value_int
		}
	}
	return stockCodes, syl1y, syl3y, syl6y
}
