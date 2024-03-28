package dto

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}
