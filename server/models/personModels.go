package models

import (
	"time"

	"github.com/jinzhu/gorm"
)
type Person struct {
	gorm.Model

	Name string
	Email string `gorm:"typevarchar(100);"`
	Birthday time.Time `gorm:"type:date;"`
	BirthdayFormatted string
	Age int
	ImagePath string
	Area string `gorm:"typevarchar(3)"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title string
	Author string
	CallNumber int `gorm:"unique_index"`
	PersonID int
}
