package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Borrow struct {
	ID         string             `bson:"_id,omitempty" json:"id,omitempty"`                // строковый ID
	ClientID   primitive.ObjectID `bson:"clientId" json:"clientId"`                         // ObjectID читателя
	BookID     primitive.ObjectID `bson:"bookId" json:"bookId"`                             // ObjectID книги
	BorrowedAt time.Time          `bson:"borrowedAt" json:"borrowedAt"`                     // Дата выдачи
	ReturnedAt *time.Time         `bson:"returnedAt,omitempty" json:"returnedAt,omitempty"` // null, если ещё не вернули
}

type BorrowStat struct {
	Date          string `bson:"date" json:"date"`                   // YYYY-MM-DD
	UniqueReaders int    `bson:"uniqueReaders" json:"uniqueReaders"` // кол-во уникальных читателей
}
