package model

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:data`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:data`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
