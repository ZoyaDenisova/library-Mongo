package domain

type Book struct {
	ID     string `bson:"_id,omitempty" json:"id,omitempty"` // строковый ID
	Title  string `bson:"title" json:"title"`                // название книги
	Author string `bson:"author" json:"author"`              // автор
	Year   int    `bson:"year" json:"year"`                  // год издания
	Genre  string `bson:"genre" json:"genre"`                // жанр
}

type BookFilter struct {
	Title  string   `json:"title"`  // фильтр по названию (нечувствительный к регистру)
	Author string   `json:"author"` // фильтр по автору
	Genres []string `json:"genres"` // один или несколько жанров
}
