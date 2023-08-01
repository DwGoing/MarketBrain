package shared

type Request struct{}

type Response struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
