package dto

import (
	"time"
)

type BorrowHistoryItem struct {
	BorrowID   string     `json:"borrowId"`
	BookID     string     `json:"bookId"`
	Title      string     `json:"title"`
	Author     string     `json:"author"`
	BorrowedAt time.Time  `json:"borrowedAt"`
	ReturnedAt *time.Time `json:"returnedAt,omitempty"`
	Status     string     `json:"status"` // "ok" / "overdue"
}

type BorrowHistoryResponse struct {
	UserID   string              `json:"userId"`
	FullName string              `json:"fullName"`
	Phone    string              `json:"phone"`
	History  []BorrowHistoryItem `json:"history"`
}

type BorrowBookInput struct {
	UserID string `json:"userId"`
	BookID string `json:"bookId"`
}

type ReturnBookInput struct {
	BorrowID string `json:"borrowId"` // id конкретной выдачи
}

type OverdueReportItem struct {
	UserID       string    `json:"userId"`
	FullName     string    `json:"fullName"`
	Phone        string    `json:"phone"`
	BookID       string    `json:"bookId"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	BorrowedAt   time.Time `json:"borrowedAt"`
	DaysOverdue  int       `json:"daysOverdue"`
	TotalOverdue int       `json:"totalOverdue"` // для повторяющихся читателей
}
