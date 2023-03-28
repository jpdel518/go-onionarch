package model

import "time"

type Article struct {
	ID            int64
	Title         string
	Content       string
	Recommend     string
	Author        Author
	UpdatedAt     time.Time
	CreatedAt     time.Time
	RecommendedAt time.Time
}
