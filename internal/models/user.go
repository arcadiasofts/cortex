package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Nickname string `gorm:"type:varchar(20);not null" json:"nickname"`
	Did      string `gorm:"type:" json:"did"`
}
