package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Status string `json:"status"`
}

type CountResponse struct {
	Count int64 `json:"count"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
