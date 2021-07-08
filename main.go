package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"richman/richman/bar"
	"strings"
)

var (
	db *gorm.DB
)

type FundStock struct {
	Id        int     `json:"id"`
	StockCode string  `json:"stock_code"`
	FundCode  string  `json:"fund_code"`
	FundName  string  `json:"fund_name"`
	Syl1y     float64 `json:"syl_1_y"`
	Syl3y     float64 `json:"syl_3_y"`
	Syl6y     float64 `json:"syl_6_y"`
}

func main() {
	setup()
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
			db.Table("fund_stock").Create(&FundStock{
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
	res, _ := Get("http://fund.eastmoney.com/js/fundcode_search.js")
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
	res, _ := Get(fmt.Sprintf("http://fund.eastmoney.com/pingzhongdata/"+"%s", code+".js"))
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

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	resBody, _ := ioutil.ReadAll(resp.Body)
	return resBody, nil
}

func setup() {
	var err error
	var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s",
		"root",
		"12345678",
		"127.0.0.1:3306",
		"fund")
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("mysql Connect Error:", err)
		db = nil
		return
	}
	db.LogMode(false)
}
