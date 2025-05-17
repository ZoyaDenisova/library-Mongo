package dto

type CreateBookInput struct {
	Title  string
	Author string
	Year   int
	Genre  string
}

type UpdateBookInput struct {
	ID     string
	Title  *string
	Author *string
	Year   *int
	Genre  *string
}
