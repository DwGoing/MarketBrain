package enum

import "fmt"

type ApiErrorType int64

const (
	ApiErrorType_Ok               ApiErrorType = 200
	ApiErrorType_RequestBindError ApiErrorType = 501
	ApiErrorType_ParameterError   ApiErrorType = 502
	ApiErrorType_ServiceError     ApiErrorType = 503
)

func (e ApiErrorType) Code() int64 {
	return int64(e)
}

func (e ApiErrorType) String(err error) string {
	var message string
	switch e {
	case ApiErrorType_Ok:
		return "success"
	case ApiErrorType_RequestBindError:
		message = "request bind error"
	case ApiErrorType_ParameterError:
		message = "parameter error"
	case ApiErrorType_ServiceError:
		message = "service error"
	}
	if err != nil {
		message = fmt.Sprintf("%s: %s", message, err)
	}
	return message
}
