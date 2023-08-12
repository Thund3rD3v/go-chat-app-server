package routes

// Responses

type ErrorResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Ok     bool `json:"Ok"`
	Uptime int  `json:"Uptime"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

// Requests

type SignUpRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
