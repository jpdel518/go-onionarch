package model

type Article struct {
	ID            int64
	Title         string
	Content       string
	Recommend     string
	AuthorId      int64
	Author        Author
	UpdatedAt     MyTime
	CreatedAt     MyTime
	RecommendedAt MyTime
}
