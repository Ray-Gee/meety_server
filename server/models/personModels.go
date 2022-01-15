package models

import (
	"time"

	"github.com/jinzhu/gorm"
)
type Person struct {
	gorm.Model

	Name string
	Email string `gorm:"typevarchar(100);unique_index"`
	Birthday time.Time `gorm:"type:date;"`
	Age int
	ImagePath string
	Books []Book
}

type Book struct {
	gorm.Model

	Title string
	Author string
	CallNumber int `gorm:"unique_index"`
	PersonID int
}
