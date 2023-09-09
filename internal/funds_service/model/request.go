package model

type Request struct{}

type Response struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
