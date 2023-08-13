package routes

import "github.com/thund3rd3v/chat-app/structs"

// Responses

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Ok     bool `json:"Ok"`
	Uptime int  `json:"Uptime"`
}

type SignInResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
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

// Websocket

type ChatEvent struct {
	EventType string `json:"eventType"`
	Value     string `json:"value"`
}

type ChatGlobalEvent struct {
	EventType string          `json:"eventType"`
	Message   structs.Message `json:"message"`
}
