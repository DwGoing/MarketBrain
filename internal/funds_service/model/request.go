package model

type Request struct {
	Id string `json:"id"`
}

type Response struct {
	Id      string `json:"id"`
	Code    int64  `json:"code"`
	Message string `json:"message,omitempty"`
}
