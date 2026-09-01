package models


import (
	"time"
	"gorm.io/gorm"
)

type Post struct{
	ID uint `gorm:"primaryKey" json:"id"`
	Title string `gorm:"size:200; not null" json:"title"`
	Content string `gorm:"type:text; not null" json:"content"`
	Author string `gorm:"size:100; not null" json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}