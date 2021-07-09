package domain

import "time"

type FundStock struct {
	Id        int       `json:"id" gorm:"primary_key"`
	StockCode string    `json:"stock_code" gorm:"type:varchar(255);DEFAULT:'';not null"`
	FundCode  string    `json:"fund_code" gorm:"type:varchar(255);DEFAULT:'';not null"`
	FundName  string    `json:"fund_name" gorm:"type:varchar(255);DEFAULT:'';not null"`
	Syl1y     float64   `json:"syl_1_y" gorm:"type:float;DEFAULT:0;not null"`
	Syl3y     float64   `json:"syl_3_y" gorm:"type:float;DEFAULT:0;not null"`
	Syl6y     float64   `json:"syl_6_y" gorm:"type:float;DEFAULT:0;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp(6)"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp(6)"`
}
