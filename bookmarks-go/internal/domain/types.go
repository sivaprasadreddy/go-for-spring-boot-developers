package domain

import "time"

type Bookmark struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"column:title"`
	Url       string    `json:"url" gorm:"column:url"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}
