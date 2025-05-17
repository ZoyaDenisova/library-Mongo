package dto

type RegisterUserInput struct {
	FullName string
	Phone    string
	Password string
	Role     string // "reader", "librarian", "admin"
}

type LoginInput struct {
	Phone    string
	Password string
}

type UpdateUserInput struct {
	ID       string
	FullName *string
	Phone    *string
	Password *string
	Role     *string
	IsActive *bool
}
