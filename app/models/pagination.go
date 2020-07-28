package models

type Pagination struct {
	Limit        int
	Page         int
	Sort         string
	TotalRows    int
	FirstPage    int
	PreviousPage int
	NextPage     string
	LastPage     string
	FromRow      string
	ToRow        int
	Rows         interface{}
}
